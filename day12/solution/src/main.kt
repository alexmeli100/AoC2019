import java.io.File
import kotlin.math.*

class Velocity(private var coords: IntArray): Iterable<Int>{
    operator fun get(i: Int) = coords[i]
    operator fun set(i: Int, value: Int) {
        coords[i] = value
    }

    override fun iterator() = coords.iterator()
    override fun toString(): String = String.format("vel=<x=${coords[0]}, y=${coords[1]}, z=${coords[2]}>")
}

class Pos(private var axis: IntArray): Iterable<Int> {
    operator fun plus(v: Velocity) {
        axis.zip(v).map { (x1, y1) -> x1+y1 }.forEachIndexed{ i, a -> axis[i] = a }
    }

    override fun iterator() = axis.iterator()
    override fun toString() = String.format("pos=<x=${axis[0]}, y=${axis[1]}, z=${axis[2]}>")

}

class Moon(private var pos: Pos, private var vel: Velocity) {
    fun applyGravity(other: Moon) {
        pos.zip(other.pos).forEachIndexed{ i, (x1, x2) ->
            val (a, b) = when(x1.compareTo(x2)) {
                1 -> Pair(-1, 1)
                -1 -> Pair(1, -1)
                else -> Pair(0, 0)
            }

            vel[i] += a
            other.vel[i] += b
        }
    }

    fun updatePos() = pos + vel

    private fun potentionEnergy() = pos.map { a -> abs(a) }.sum()
    private fun kineticEnergy() = vel.map{ a -> abs(a)}.sum()
    fun totalEnergy() = potentionEnergy() * kineticEnergy()

    override fun toString() = String.format("$pos  $vel")
}

class Galaxy(private var moons: MutableList<Moon>) {
    private fun updateVel() {
        for (i in 0 until moons.size) {
            for (j in i+1 until moons.size) {
                moons[i].applyGravity(moons[j])
            }
        }
    }

    fun update() {
        updateVel()
        moons.forEach { it.updatePos() }
    }

    fun moonEnergies() = moons.map{ it.totalEnergy() }.sum()
}

private val input by lazy { File("input.txt").readText() }

fun main() {
    val pat = "[-]?\\d+".toRegex()
    val moons = mutableListOf<Moon>()

    input.lineSequence().forEach { s ->
        val pos = IntArray(3)
        pat.findAll(s).map { it.value.toInt() }.forEachIndexed { i, x -> pos[i] = x }
        moons.add(Moon(Pos(pos), Velocity(intArrayOf(0, 0, 0))))
    }

    val galaxy = Galaxy(moons)

    for (i in 0 until 1000)
        galaxy.update()

    println("Total Energy is ${galaxy.moonEnergies()}")
}