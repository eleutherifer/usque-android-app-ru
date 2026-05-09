package com.diniboy.usqueandroid

import android.app.Service
import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor
import android.system.Os
import android.util.Log
import usqueandroid.Usqueandroid
import usqueandroid.VpnStateCallback
import java.io.File
import java.io.FileDescriptor
import java.util.concurrent.Executors
import java.util.concurrent.atomic.AtomicBoolean

class UsqueVpnService : VpnService() {
    companion object {
        const val ACTION_STOP = "com.diniboy.usqueandroid.STOP_VPN"
        private const val TAG = "UsqueVpnService"
        @Volatile private var activeService: UsqueVpnService? = null

        fun stopActiveTunnel() {
            activeService?.stopVpn("external stop") ?: runCatching { Usqueandroid.stopTunnel() }
        }
    }

    private val executor = Executors.newSingleThreadExecutor()
    private var tun: ParcelFileDescriptor? = null
    private var detachedTunFd: Int = -1
    private val running = AtomicBoolean(false)

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        activeService = this
        if (intent?.action == ACTION_STOP) {
            Log.i(TAG, "stop requested")
            stopVpn("ACTION_STOP")
            stopSelf()
            return Service.START_NOT_STICKY
        }

        val configPath = intent?.getStringExtra("configPath") ?: File(filesDir, "config.json").absolutePath
        val sni = intent?.getStringExtra("sni") ?: "www.visa.cn"
        val endpoint = intent?.getStringExtra("endpoint") ?: "162.159.198.2:500"
        val splitMode = intent?.getBooleanExtra("splitMode", false) ?: false
        val allowedApps = intent?.getStringArrayListExtra("allowedApps") ?: arrayListOf()

        if (running.get()) return Service.START_STICKY
        running.set(true)

        executor.execute {
            try {
                Log.i(TAG, "starting vpn service endpoint=$endpoint sni=$sni splitMode=$splitMode allowedApps=${allowedApps.size} config=$configPath")
                Usqueandroid.resetConnectionOptions()
                Usqueandroid.setSNI(sni)
                Usqueandroid.setEndpoint(endpoint)
                Log.i(TAG, "native endpoint now=${runCatching { Usqueandroid.getEndpoint() }.getOrDefault("")}")

                val builder = Builder()
                    .setSession("Usque VPN")
                    .setMtu(1280)
                    .addAddress(safeIPv4(configPath), 32)
                    .addDnsServer("1.1.1.1")
                    .addDnsServer("1.0.0.1")
                    .addRoute("0.0.0.0", 0)

                if (splitMode) {
                    if (allowedApps.isEmpty()) throw IllegalStateException("split mode enabled but no apps selected")
                    allowedApps.distinct().forEach { pkg ->
                        runCatching { builder.addAllowedApplication(pkg) }
                            .onFailure { Log.w(TAG, "addAllowedApplication failed: $pkg", it) }
                    }
                } else {
                    // Critical: do not route this app's own MASQUE/QUIC control connection into itself.
                    runCatching { builder.addDisallowedApplication(packageName) }
                        .onFailure { Log.w(TAG, "addDisallowedApplication failed", it) }
                }

                val ipv6 = runCatching { Usqueandroid.getAssignedIPv6(configPath) }.getOrDefault("")
                if (ipv6.isNotBlank()) runCatching {
                    builder.addAddress(ipv6, 128)
                    builder.addRoute("::", 0)
                    builder.addDnsServer("2606:4700:4700::1111")
                    builder.addDnsServer("2606:4700:4700::1001")
                }.onFailure { Log.w(TAG, "ipv6 setup failed", it) }

                val pfd = builder.establish() ?: throw IllegalStateException("builder.establish returned null")
                tun = pfd
                detachedTunFd = pfd.detachFd()
                tun = null
                Log.i(TAG, "tun established fd=$detachedTunFd")

                // Native fd mode: Go owns the Android TUN fd and handles the full data plane.
                // Do NOT pass connect-port here. The second argument is tunFd.
                val err = Usqueandroid.startTunnelWithFd(configPath, detachedTunFd.toLong(), object : VpnStateCallback {
                    override fun onConnected() { Log.i(TAG, "tunnel connected") }
                    override fun onDisconnected(reason: String?) {
                        Log.w(TAG, "tunnel disconnected: $reason")
                        stopVpn("native disconnected")
                        stopSelf()
                    }
                    override fun onError(message: String?) {
                        Log.e(TAG, "tunnel error: $message")
                        stopVpn("native error")
                        stopSelf()
                    }
                })
                if (!err.isNullOrBlank()) throw IllegalStateException(err)

                Log.i(TAG, "startTunnelWithFd returned without error")
            } catch (e: Exception) {
                Log.e(TAG, "vpn service failed", e)
                stopVpn("exception")
                stopSelf()
            }
        }
        return Service.START_STICKY
    }

    private fun safeIPv4(configPath: String): String {
        return runCatching { Usqueandroid.getAssignedIPv4(configPath) }
            .getOrDefault("")
            .ifBlank { "172.16.0.2" }
    }

    private fun stopVpn(reason: String = "stop") {
        Log.i(TAG, "stopping vpn: $reason fd=$detachedTunFd running=${running.get()}")
        running.set(false)
        runCatching { Usqueandroid.stopTunnel() }
            .onFailure { Log.w(TAG, "native stopTunnel failed", it) }
        runCatching { tun?.close() }
            .onFailure { Log.w(TAG, "tun close failed", it) }
        tun = null
        if (detachedTunFd >= 0) {
            runCatching { Os.close(fileDescriptorFromInt(detachedTunFd)) }
                .onFailure { Log.w(TAG, "detached fd close failed", it) }
            detachedTunFd = -1
        }
    }

    private fun fileDescriptorFromInt(fdInt: Int): FileDescriptor {
        val fd = FileDescriptor()
        val field = FileDescriptor::class.java.getDeclaredField("descriptor")
        field.isAccessible = true
        field.setInt(fd, fdInt)
        return fd
    }

    override fun onDestroy() {
        stopVpn("onDestroy")
        if (activeService === this) activeService = null
        super.onDestroy()
    }
}
