package com.warp.usque

import java.text.SimpleDateFormat
import java.util.ArrayDeque
import java.util.Date
import java.util.Locale

/**
 * Простой кольцевой буфер логов, который живёт прямо в приложении —
 * чтобы можно было прочитать и скопировать диагностику без adb и компьютера.
 * Ничего не пишет на диск, всё только в памяти процесса.
 */
object DiagLog {
    private const val MAX_LINES = 500
    private val lines = ArrayDeque<String>()
    private val fmt = SimpleDateFormat("HH:mm:ss.SSS", Locale.US)

    @Synchronized
    fun add(tag: String, msg: String) {
        lines.addLast("${fmt.format(Date())} [$tag] $msg")
        while (lines.size > MAX_LINES) lines.removeFirst()
    }

    @Synchronized
    fun getAll(): String = lines.joinToString("\n")

    @Synchronized
    fun clear() = lines.clear()
}
