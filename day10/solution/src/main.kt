import java.io.File
import kotlin.math.*

data class Pos(val x: Int, val y: Int) {
    operator fun minus(other: Pos) = Pos(x - other.x, y - other.y)
    operator fun times(scale: Int) = Pos(x * scale, y * scale)
    infix fun cross(b: Pos) = x.toLong() * b.y - y.toLong() * b.x
    fun manDist(other: Pos) = abs(x - other.x) + abs(y - other.y)
    fun manDist() = abs(x) + abs(y)

    fun normalize() = gcd(x, y).let { Pos(x/it, y/it) }


}

private val input by lazy { File("input.txt").readText() }

fun main() {
    val asteriods = mutableListOf<Pos>()
    input.lineSequence().forEachIndexed { y, ln ->
        ln.forEachIndexed { x, c ->
            if (c == '#') asteriods.add(Pos(x, y))
        }
    }

    val (best, ans1) = asteriods.map { a ->
        val hits = HashSet<Pos>()
        for (b in asteriods) {
            if (a != b) hits.add((b - a).normalize())
        }
        a to hits.size
    }.maxBy { it.second }!!

    println("Part 1: $ans1")

}

tailrec fun gcd(a: Int, b: Int): Int = if(a == 0) abs(b) else gcd(b % a, a)