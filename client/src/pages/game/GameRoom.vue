<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import * as THREE from 'three'
import TavernShell from '@/components/wrappers/TavernShell.vue'
import UiButton from '@/components/form/UiButton.vue'

const route = useRoute()
const roomCode = computed(() => String(route.params.code ?? '').toUpperCase())
const pointsToScore = 3000

// UI / flow
const hasStarted = ref(false)

// DOM
const canvasEl = ref<HTMLCanvasElement | null>(null)
const wrapEl = ref<HTMLDivElement | null>(null)

// three
let renderer: THREE.WebGLRenderer | null = null
let scene: THREE.Scene | null = null
let camera: THREE.PerspectiveCamera | null = null
let raf = 0

let raycaster: THREE.Raycaster | null = null
const pointer = new THREE.Vector2()

// avoid allocations in render loop
const tmpV3 = new THREE.Vector3()

// camera base + push-in + shake
const camBasePos = new THREE.Vector3()
const camBaseLook = new THREE.Vector3()
const camRollOffset = new THREE.Vector3(0, -0.35, -0.45)
let camZoomT = 0
const shake = { amp: 0, seed: Math.random() * 1000 }

// sizes (compact dice)
const D = 0.62
const D_HALF = D / 2
const DIE_Y = D_HALF + 0.02

// shared assets
let diceGeo: THREE.BoxGeometry | null = null
let diceMats: THREE.MeshStandardMaterial[] | null = null

let tableGeo: THREE.PlaneGeometry | null = null
let tableMat: THREE.MeshStandardMaterial | null = null
let tableMesh: THREE.Mesh | null = null

// cup
let cup: THREE.Mesh | null = null
let cupGeo: THREE.BufferGeometry | null = null
let cupMat: THREE.MeshStandardMaterial | null = null
const cupBaseRot = new THREE.Euler()

// cinematic cup positions
const cupOffPos = new THREE.Vector3(-4.2, 0.2, 0.9) // off-left
const cupCenterPos = new THREE.Vector3(-0.9, 0.25, 0.55) // centered shot
const cupLiftOutPos = new THREE.Vector3(0.3, 2.2, 2.6) // up & away

type Phase = 'idle' | 'cup_enter' | 'cup_shake' | 'cup_place' | 'cup_lift_out' | 'dice_settle'
const phase = ref<Phase>('idle')

// timings
let phaseStartAt = 0
const CUP_ENTER_MS = 220
const CUP_SHAKE_MS = 520
const CUP_PLACE_MS = 220
const CUP_LIFT_MS = 260
const DICE_DROP_MS = 420
const DICE_SETTLE_MS = 320

type Die = {
  id: string
  mesh: THREE.Mesh
  value: number

  inPlay: boolean
  selected: boolean

  playPos: THREE.Vector3

  dropFrom: THREE.Vector3
  dropTo: THREE.Vector3

  settling: boolean
  settleDelayMs: number
  settleMs: number
  startQuat: THREE.Quaternion
  endQuat: THREE.Quaternion
}

const dice = ref<Die[]>([])

// turn state (KCD-like)
const totalScore = ref(0)
const turnScore = ref(0)
const lastThrowScorePreview = ref<number | null>(null)
const message = ref<string>('')

// rule gate: after a scoring throw you must take at least one scoring die before rolling again
const mustTakeAfterThrow = ref(false)

// derived
const rolling = computed(() => phase.value !== 'idle')
const diceInPlayCount = computed(() => dice.value.filter((d) => d.inPlay).length)
const selectedCount = computed(() => dice.value.filter((d) => d.inPlay && d.selected).length)
const resultsInPlay = computed(() => dice.value.filter((d) => d.inPlay).map((d) => d.value))

const canTake = computed(() => {
  if (rolling.value) return false
  if (!selectedCount.value) return false
  return evaluateSelection(selectedValues()).valid
})

const canRollNow = computed(() => {
  if (rolling.value) return false
  if (message.value.startsWith('FARKLE')) return false
  if (mustTakeAfterThrow.value && selectedCount.value === 0) return false
  return true
})

const canBank = computed(() => canTake.value)

// ---------- landing area (spread dice on table) ----------
const LAND_CENTER = new THREE.Vector3(-0.15, DIE_Y, 0.25)
const LAND_RX = 1.55
const LAND_RZ = 1.05
const MIN_DIST = 0.72

function randInEllipse(cx: number, cz: number, rx: number, rz: number) {
  const a = Math.random() * Math.PI * 2
  const r = Math.sqrt(Math.random())
  return { x: cx + Math.cos(a) * rx * r, z: cz + Math.sin(a) * rz * r }
}

function generateNonOverlappingTargets(count: number) {
  const out: Array<{ x: number; z: number }> = []
  const MAX_TRIES = 90

  for (let i = 0; i < count; i++) {
    let chosen = { x: LAND_CENTER.x, z: LAND_CENTER.z }

    for (let tries = 0; tries < MAX_TRIES; tries++) {
      const p = randInEllipse(LAND_CENTER.x, LAND_CENTER.z, LAND_RX, LAND_RZ)

      let good = true
      for (const q of out) {
        const dx = p.x - q.x
        const dz = p.z - q.z
        if (dx * dx + dz * dz < MIN_DIST * MIN_DIST) {
          good = false
          break
        }
      }
      if (good) {
        chosen = p
        break
      }
    }

    out.push(chosen)
  }

  return out
}

// ---------- helpers ----------
const clamp01 = (x: number) => Math.min(1, Math.max(0, x))
const easeOutCubic = (t: number) => 1 - Math.pow(1 - t, 3)
const easeInOutCubic = (t: number) =>
  t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2

// IMPORTANT: shake should ADD to current camera position (not overwrite base),
// otherwise it cancels our zoom/push-in.
function applyCameraShake(t: number) {
  if (!camera) return
  const a = shake.amp
  if (a <= 0.0001) return

  const s = t * 0.01 + shake.seed
  const ox = (Math.sin(s * 9.1) + Math.sin(s * 13.7) * 0.5) * a
  const oy = (Math.cos(s * 11.3) + Math.sin(s * 7.9) * 0.5) * a * 0.6
  const oz = (Math.sin(s * 8.2) + Math.cos(s * 10.4) * 0.5) * a * 0.8

  camera.position.x += ox
  camera.position.y += oy
  camera.position.z += oz

  shake.amp *= 0.92
}

// ---------- scoring ----------
function countFaces(values: number[]) {
  const c = [0, 0, 0, 0, 0, 0, 0]
  for (const v of values) {
    if (c[v]) {
      c[v]++
    }
  }
  return c
}

function isStraight(values: number[], target: number[]) {
  if (values.length !== target.length) return false
  const sorted = [...values].sort((a, b) => a - b)
  for (let i = 0; i < target.length; i++) if (sorted[i] !== target[i]) return false
  return true
}

function hasAnyScore(values: number[]) {
  if (!values.length) return false
  const c = countFaces(values)
  if (c[1] && c[1] > 0 || c[5] && c[5] > 0) return true
  for (let f = 1; f <= 6; f++) if ((c[f] ?? 0) >= 3) return true
  return (
    isStraight(values, [1, 2, 3, 4, 5, 6]) ||
    isStraight(values, [1, 2, 3, 4, 5]) ||
    isStraight(values, [2, 3, 4, 5, 6])
  )
}

function evaluateSelection(values: number[]): { valid: boolean; score: number; label: string } {
  if (!values.length) return { valid: false, score: 0, label: 'Select dice' }

  if (isStraight(values, [1, 2, 3, 4, 5, 6])) return { valid: true, score: 1500, label: 'Straight 1–6' }
  if (isStraight(values, [1, 2, 3, 4, 5])) return { valid: true, score: 500, label: 'Straight 1–5' }
  if (isStraight(values, [2, 3, 4, 5, 6])) return { valid: true, score: 750, label: 'Straight 2–6' }

  const c = countFaces(values)
  let score = 0

  for (let f = 1; f <= 6; f++) {
    const n = c[f]
    if (!n) continue

    if (n >= 3) {
      const base = f === 1 ? 1000 : f * 100
      const mult = Math.pow(2, n - 3)
      score += base * mult
      continue
    }

    if (f === 1) score += n * 100
    else if (f === 5) score += n * 50
    else return { valid: false, score: 0, label: 'Invalid selection' }
  }

  return { valid: score > 0, score, label: score > 0 ? 'Scoring dice' : 'Invalid selection' }
}

function selectedValues() {
  return dice.value.filter((d) => d.inPlay && d.selected).map((d) => d.value)
}

// ---------- textures ----------
function makePipTexture(value: number) {
  const size = 256
  const c = document.createElement('canvas')
  c.width = size
  c.height = size
  const ctx = c.getContext('2d')!

  ctx.fillStyle = '#F2E6CF'
  ctx.fillRect(0, 0, size, size)

  ctx.strokeStyle = 'rgba(60,40,25,0.22)'
  ctx.lineWidth = 10
  ctx.strokeRect(8, 8, size - 16, size - 16)

  const pip = (x: number, y: number) => {
    ctx.beginPath()
    ctx.arc(x, y, 18, 0, Math.PI * 2)
    ctx.fill()
  }
  ctx.fillStyle = 'rgba(35,20,12,0.92)'

  const pad = 64
  const mid = size / 2
  const left = pad
  const right = size - pad
  const top = pad
  const bottom = size - pad

  const layouts: Record<number, [number, number][]> = {
    1: [[mid, mid]],
    2: [[left, top], [right, bottom]],
    3: [[left, top], [mid, mid], [right, bottom]],
    4: [[left, top], [right, top], [left, bottom], [right, bottom]],
    5: [[left, top], [right, top], [mid, mid], [left, bottom], [right, bottom]],
    6: [[left, top], [left, mid], [left, bottom], [right, top], [right, mid], [right, bottom]],
  }
  layouts[value]?.forEach(([x, y]) => pip(x, y))

  const tex = new THREE.CanvasTexture(c)
  tex.colorSpace = THREE.SRGBColorSpace
  tex.anisotropy = 4
  tex.needsUpdate = true
  return tex
}

function makeDiceMaterials() {
  // BoxGeometry material order: [+X, -X, +Y, -Y, +Z, -Z]
  // mapping: +Y=1, -Y=6, +Z=2, -Z=5, +X=3, -X=4
  const tex1 = makePipTexture(1)
  const tex2 = makePipTexture(2)
  const tex3 = makePipTexture(3)
  const tex4 = makePipTexture(4)
  const tex5 = makePipTexture(5)
  const tex6 = makePipTexture(6)

  const mat = (tex: THREE.Texture) =>
    new THREE.MeshStandardMaterial({
      map: tex,
      roughness: 0.65,
      metalness: 0.05,
      emissive: new THREE.Color(0x000000),
      emissiveIntensity: 0,
    })

  return [mat(tex3), mat(tex4), mat(tex1), mat(tex6), mat(tex2), mat(tex5)]
}

function quatForTop(value: number) {
  const e = new THREE.Euler(0, 0, 0, 'XYZ')
  if (value === 1) e.set(0, 0, 0)
  if (value === 6) e.set(Math.PI, 0, 0)
  if (value === 2) e.set(-Math.PI / 2, 0, 0)
  if (value === 5) e.set(Math.PI / 2, 0, 0)
  if (value === 3) e.set(0, 0, Math.PI / 2)
  if (value === 4) e.set(0, 0, -Math.PI / 2)

  const q = new THREE.Quaternion()
  q.setFromEuler(e)
  return q
}

// ---------- visuals ----------
function applyDieVisual(d: Die) {
  if (!Array.isArray(d.mesh.material)) return

  // якщо кубик забрали і продовжили хід
  if (!d.inPlay) {
    d.mesh.visible = false

    return
  }

  d.mesh.scale.set(1, 1, 1)
  // d.mesh.material.forEach((m) => {
    // m.emissive.setHex(d.selected ? 0x6b4a1b : 0x000000)
    // m.emissiveIntensity = d.selected ? 0.14 : 0
  // })
}

function createDie(playX: number, playZ: number, id: string) {
  if (!scene || !diceGeo || !diceMats) throw new Error('three not ready')

  const mesh = new THREE.Mesh(diceGeo, diceMats)
  mesh.castShadow = true
  mesh.position.set(playX, DIE_Y, playZ)
  mesh.quaternion.copy(quatForTop(1))
  mesh.rotation.setFromQuaternion(mesh.quaternion)
  mesh.userData = { dieId: id }

  const d: Die = {
    id,
    mesh,
    value: 1,

    inPlay: true,
    selected: false,

    playPos: new THREE.Vector3(playX, DIE_Y, playZ),

    dropFrom: new THREE.Vector3(),
    dropTo: new THREE.Vector3(playX, DIE_Y, playZ),

    settling: false,
    settleDelayMs: 0,
    settleMs: DICE_SETTLE_MS,
    startQuat: new THREE.Quaternion(),
    endQuat: new THREE.Quaternion(),
  }

  scene.add(mesh)
  applyDieVisual(d)
  return d
}

// ---------- cinematic roll ----------
function prepareThrowTargets() {
  if (!cup) return
  lastThrowScorePreview.value = null

  const inPlay = dice.value.filter((d) => d.inPlay)
  const targets = generateNonOverlappingTargets(inPlay.length)

  inPlay.forEach((d, idx) => {
    d.selected = false
    applyDieVisual(d)

    d.value = Math.floor(Math.random() * 6) + 1

    d.dropFrom.set(cupCenterPos.x + 0.22, cupCenterPos.y + 0.82, cupCenterPos.z - 0.1)

    if (targets[idx]) {
      d.dropTo.set(targets[idx].x, DIE_Y, targets[idx].z)
    }

    d.settling = true
    d.settleDelayMs = idx * 55
    d.settleMs = DICE_SETTLE_MS + Math.random() * 120
    d.startQuat.copy(d.mesh.quaternion)

    const target = quatForTop(d.value)
    const yaw = new THREE.Quaternion().setFromAxisAngle(new THREE.Vector3(0, 1, 0), Math.random() * Math.PI * 2)
    d.endQuat.copy(yaw.multiply(target))
  })
}

function startCinematicRoll() {
  if (!cup) return
  if (phase.value !== 'idle') return
  if (diceInPlayCount.value === 0) return

  camZoomT = 0

  if (!hasStarted.value) {
    hasStarted.value = true
    message.value = ''
    if (tableMesh) tableMesh.visible = true
  }

  message.value = ''
  prepareThrowTargets()

  cup.visible = true
  cup.position.copy(cupOffPos)
  cup.rotation.copy(cupBaseRot)

  for (const d of dice.value) {
    if (d.inPlay) d.mesh.visible = false
  }

  phase.value = 'cup_enter'
  phaseStartAt = performance.now()
}

function finalizeAfterThrow() {
  const vals = resultsInPlay.value
  if (!hasAnyScore(vals)) {
    turnScore.value = 0
    message.value = 'FARKLE! No scoring dice. You lost the points of this turn.'
    lastThrowScorePreview.value = 0
    mustTakeAfterThrow.value = false
    return
  }

  mustTakeAfterThrow.value = true
  message.value = 'Select scoring dice, then Take & Roll or Take & Bank.'
  lastThrowScorePreview.value = null
}

// ---------- gameplay actions ----------
function toggleSelectDieById(id: string) {
  if (phase.value !== 'idle') return
  if (message.value.startsWith('FARKLE')) return

  const d = dice.value.find((x) => x.id === id)
  if (!d || !d.inPlay) return

  d.selected = !d.selected
  applyDieVisual(d)

  const evalRes = evaluateSelection(selectedValues())
  lastThrowScorePreview.value = evalRes.valid ? evalRes.score : null
}

function takeSelectedAndContinue() {
  const values = selectedValues()
  const evalRes = evaluateSelection(values)
  if (!evalRes.valid) return

  turnScore.value += evalRes.score
  message.value = `+${evalRes.score} (${evalRes.label}) — rolling…`
  lastThrowScorePreview.value = null
  mustTakeAfterThrow.value = false

  dice.value.forEach((d) => {
    if (d.inPlay && d.selected) {
      d.inPlay = false
      d.selected = false
      applyDieVisual(d)
    }
  })

  // hot dice: put all dice back in play, but DO NOT snap to grid (looks wrong now).
  if (diceInPlayCount.value === 0) {
    message.value = `HOT DICE! (+${evalRes.score}). All dice back in play.`
    const targets = generateNonOverlappingTargets(dice.value.length)

    dice.value.forEach((d, idx) => {
      d.inPlay = true
      d.selected = false
      d.mesh.visible = true
      if (targets[idx]) {
        d.mesh.position.set(targets[idx].x, DIE_Y, targets[idx].z)
      }
      d.dropTo.copy(d.mesh.position)
      applyDieVisual(d)
    })
  }

  startCinematicRoll()
}

function takeSelectedAndBank() {
  const values = selectedValues()
  const evalRes = evaluateSelection(values)
  if (!evalRes.valid) return

  turnScore.value += evalRes.score
  totalScore.value += turnScore.value
  message.value = `Banked ${turnScore.value} points.`
  lastThrowScorePreview.value = null
  mustTakeAfterThrow.value = false
  resetTurn()
}

function resetTurn() {
  const targets = generateNonOverlappingTargets(dice.value.length)

  dice.value.forEach((d, idx) => {
    d.inPlay = true
    d.selected = false
    d.mesh.visible = true

    if (targets[idx]) {
      d.mesh.position.set(targets[idx].x, DIE_Y, targets[idx].z)
    }

    d.dropTo.copy(d.mesh.position)
    applyDieVisual(d)
  })

  turnScore.value = 0
}

function roll() {
  if (!canRollNow.value) {
    if (mustTakeAfterThrow.value && selectedCount.value === 0) {
      message.value = 'Take scoring dice first.'
    }
    return
  }
  startCinematicRoll()
}

// ---------- init three ----------
function initThree() {
  if (!canvasEl.value || !wrapEl.value) return

  renderer = new THREE.WebGLRenderer({ canvas: canvasEl.value, antialias: true, alpha: true })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
  renderer.outputColorSpace = THREE.SRGBColorSpace
  renderer.shadowMap.enabled = true

  scene = new THREE.Scene()
  raycaster = new THREE.Raycaster()

  camera = new THREE.PerspectiveCamera(32, 1, 0.1, 100)
  camera.position.set(0.0, 7.2, 4.6)
  camera.lookAt(0.0, -0.6, 0.35)
  camBasePos.copy(camera.position)
  camBaseLook.set(0.0, -0.6, 0.35)

  const ambient = new THREE.AmbientLight(0xffffff, 0.45)
  scene.add(ambient)

  const key = new THREE.DirectionalLight(0xffffff, 0.95)
  key.position.set(3, 6, 2)
  key.castShadow = true
  scene.add(key)

  const rim = new THREE.DirectionalLight(0xffffff, 0.35)
  rim.position.set(-4, 3, -3)
  scene.add(rim)

  // (optional) soft warm fill to make dice readable
  const fill = new THREE.PointLight(0xffe6c2, 0.75, 18, 2)
  fill.position.set(0.0, 3.2, 1.2)
  scene.add(fill)

  // table (ONLY ONE)
  tableGeo = new THREE.PlaneGeometry(14, 9)
  tableMat = new THREE.MeshStandardMaterial({ color: 0x2a1b12, roughness: 0.9, metalness: 0.05 })
  tableMesh = new THREE.Mesh(tableGeo, tableMat)
  tableMesh.rotation.x = -Math.PI / 2
  tableMesh.position.y = -0.6
  tableMesh.receiveShadow = true
  tableMesh.visible = false
  scene.add(tableMesh)

  // cup
  cupGeo = new THREE.CylinderGeometry(0.75, 0.95, 1.7, 28, 1, true)
  cupMat = new THREE.MeshStandardMaterial({
    color: 0x3a2418,
    roughness: 0.95,
    metalness: 0.02,
    side: THREE.DoubleSide,
  })
  cup = new THREE.Mesh(cupGeo, cupMat)
  cup.castShadow = true
  cup.visible = false
  cup.rotation.copy(cupBaseRot)
  scene.add(cup)

  // dice shared
  diceGeo = new THREE.BoxGeometry(D, D, D)
  diceMats = makeDiceMaterials()

  // initial positions (don’t matter much; we hide them until start)
  const play: Array<[number, number]> = [
    [-0.7, 0.55],
    [0.0, 0.55],
    [0.7, 0.55],
    [-0.7, -0.1],
    [0.0, -0.1],
    [0.7, -0.1],
  ]

  dice.value = play.map(([px, pz], i) => {
    return createDie(px, pz, `d${i + 1}`)
  })

  // hide dice before first roll
  dice.value.forEach((d) => (d.mesh.visible = false))

  const resize = () => {
    if (!renderer || !camera || !wrapEl.value) return
    const rect = wrapEl.value.getBoundingClientRect()
    const w = Math.max(1, Math.floor(rect.width))
    const h = Math.max(1, Math.floor(rect.height))
    renderer.setSize(w, h, false)
    camera.aspect = w / h
    camera.updateProjectionMatrix()
  }

  const onClick = (ev: MouseEvent) => {
    if (!raycaster || !camera || !wrapEl.value) return
    if (phase.value !== 'idle') return
    if (!hasStarted.value) return

    const rect = wrapEl.value.getBoundingClientRect()
    const x = ((ev.clientX - rect.left) / rect.width) * 2 - 1
    const y = -(((ev.clientY - rect.top) / rect.height) * 2 - 1)
    pointer.set(x, y)

    raycaster.setFromCamera(pointer, camera)
    const meshes = dice.value.map((d) => d.mesh)
    const hits = raycaster.intersectObjects(meshes, false)
    if (!hits.length) return

    const hit = hits[0]?.object as THREE.Mesh
    const id = String(hit.userData?.dieId ?? '')
    toggleSelectDieById(id)
  }

  wrapEl.value.addEventListener('click', onClick)

  const animate = (t: number) => {
    if (!renderer || !scene || !camera || !cup) return

    // hard guarantee: nothing 3D visible before first roll
    if (!hasStarted.value) {
      cup.visible = false
      if (tableMesh) tableMesh.visible = false
      for (const d of dice.value) d.mesh.visible = false
    }

    if (phase.value === 'cup_enter') {
      const s = clamp01((t - phaseStartAt) / CUP_ENTER_MS)
      const k = easeOutCubic(s)
      cup.position.lerpVectors(cupOffPos, cupCenterPos, k)
      cup.rotation.z = cupBaseRot.z + 0.06 * Math.sin(k * Math.PI)
      if (s >= 1) {
        phase.value = 'cup_shake'
        phaseStartAt = t
        shake.amp = Math.max(shake.amp, 0.06)
      }
    }

    if (phase.value === 'cup_shake') {
      const s = clamp01((t - phaseStartAt) / CUP_SHAKE_MS)
      const k = easeOutCubic(s)

      const amp = (1 - k) * 0.16
      cup.position.x = cupCenterPos.x + (Math.random() - 0.5) * amp
      cup.position.z = cupCenterPos.z + (Math.random() - 0.5) * amp
      cup.rotation.z = cupBaseRot.z + Math.sin(s * 32) * amp * 0.9
      cup.rotation.x = cupBaseRot.x + Math.cos(s * 28) * amp * 0.6

      if (s >= 1) {
        cup.position.copy(cupCenterPos)
        phase.value = 'cup_place'
        phaseStartAt = t
      }
    }

    if (phase.value === 'cup_place') {
      const s = clamp01((t - phaseStartAt) / CUP_PLACE_MS)
      const k = easeInOutCubic(s)

      cup.rotation.x = THREE.MathUtils.lerp(cup.rotation.x, cupBaseRot.x + 0.10, k)
      cup.rotation.z = THREE.MathUtils.lerp(cup.rotation.z, cupBaseRot.z, k)
      cup.position.y = THREE.MathUtils.lerp(cup.position.y, cupCenterPos.y, k)

      if (s >= 1) {
        // EMIT DICE
        for (const d of dice.value) {
          if (!d.inPlay) continue
          d.mesh.visible = true
          d.dropFrom.set(cup.position.x + 0.22, cup.position.y + 0.82, cup.position.z - 0.1)
          d.mesh.position.copy(d.dropFrom)
        }

        shake.amp = Math.max(shake.amp, 0.09)

        phase.value = 'cup_lift_out'
        phaseStartAt = t
      }
    }

    if (phase.value === 'cup_lift_out') {
      const s = clamp01((t - phaseStartAt) / CUP_LIFT_MS)
      const k = easeOutCubic(s)

      cup.position.lerpVectors(cupCenterPos, cupLiftOutPos, k)
      cup.rotation.x = cupBaseRot.x + 0.35 * k
      cup.rotation.z = cupBaseRot.z - 0.18 * k

      const sd = clamp01((t - phaseStartAt) / DICE_DROP_MS)
      const kd = easeInOutCubic(sd)

      for (const d of dice.value) {
        if (!d.inPlay) continue

        tmpV3.lerpVectors(d.dropFrom, d.dropTo, kd)
        const h = 1.05
        tmpV3.y += 4 * h * kd * (1 - kd)
        d.mesh.position.copy(tmpV3)
      }

      if (sd >= 1 && s >= 1) {
        cup.visible = false
        shake.amp = Math.max(shake.amp, 0.04)

        for (const d of dice.value) {
          if (!d.inPlay) continue
          d.mesh.position.copy(d.dropTo)
          d.mesh.position.y = DIE_Y
        }

        phase.value = 'dice_settle'
        phaseStartAt = t
      }
    }

    if (phase.value === 'dice_settle') {
      const dt = t - phaseStartAt
      let any = false

      for (const d of dice.value) {
        if (!d.inPlay || !d.settling) continue
        any = true

        const localT = dt - d.settleDelayMs
        if (localT < 0) continue

        const st = clamp01(localT / d.settleMs)
        const kk = easeOutCubic(st)
        d.mesh.quaternion.slerpQuaternions(d.startQuat, d.endQuat, kk)

        if (st >= 1) {
          d.settling = false
          d.mesh.quaternion.copy(d.endQuat)
          d.mesh.rotation.setFromQuaternion(d.mesh.quaternion)
          applyDieVisual(d)
        }
      }

      if (!any) {
        phase.value = 'idle'
        finalizeAfterThrow()
      }
    }

    // camera push-in/out
    if (phase.value !== 'idle') camZoomT = Math.min(1, camZoomT + 0.04)
    else camZoomT = Math.max(0, camZoomT - 0.04)

    tmpV3.copy(camBasePos).addScaledVector(camRollOffset, easeOutCubic(camZoomT))
    camera.position.copy(tmpV3)
    camera.lookAt(camBaseLook)

    // add shake AFTER zoom
    applyCameraShake(t)

    renderer.render(scene, camera)
    raf = requestAnimationFrame(animate)
  }

  resize()
  raf = requestAnimationFrame(animate)

  const ro = new ResizeObserver(resize)
  ro.observe(wrapEl.value)

  onBeforeUnmount(() => {
    cancelAnimationFrame(raf)
    ro.disconnect()
    wrapEl.value?.removeEventListener('click', onClick)

    dice.value.forEach((d) => scene?.remove(d.mesh))
    dice.value = []

    if (cup) scene?.remove(cup)
    if (tableMesh) scene?.remove(tableMesh)

    diceGeo?.dispose()
    diceGeo = null

    diceMats?.forEach((m) => {
      m.map?.dispose()
      m.dispose()
    })
    diceMats = null

    tableGeo?.dispose()
    tableGeo = null
    tableMat?.dispose()
    tableMat = null
    tableMesh = null

    cupGeo?.dispose()
    cupGeo = null
    cupMat?.dispose()
    cupMat = null
    cup = null

    renderer?.dispose()
    renderer = null
    scene = null
    camera = null
    raycaster = null
  })
}

onMounted(() => {
  initThree()
  message.value = '' // overlay handles the start text
})
</script>

<template>
  <TavernShell >
    <div class="grid gap-4 lg:grid-cols-12 min-h-[60dvh] lg:min-h-[65dvh] mt-8">
      <div class="lg:col-span-12">
        <div class="mt-3 text-xl text-parchment-50/60 flex justify-between flex-wrap items-center gap-2">
          <div class="flex gap-5">
             <span>
              Room: <span class="text-parchment-50 font-semibold">{{ roomCode }}</span>
            </span>
            <span>
              Bet: <span class="text-parchment-50 font-semibold">{{ 100 }} gold</span>
            </span>
          </div>
          <RouterLink to="/" >Leave room →</RouterLink>
        </div>
        <div class="mt-4 rounded-[var(--radius-tavern)] border border-wood-700/25 bg-tavern-900/35 overflow-hidden">
          <div ref="wrapEl" class="relative h-[52dvh] lg:h-[65dvh] select-none">
            <canvas ref="canvasEl" class="block h-full w-full"></canvas>

            <!-- Opponent score -->
            <div
              class="absolute top-4 left-4 flex flex-col gap-2 z-20"
            >
              <div class="lg:col-span-4">
                <div class="panel">
                  <div class="font-display text-xl">Opponent</div>
                  <div class="mt-4 rounded-xl border border-wood-700/25 bg-parchment-50/60 p-4 space-y-3">
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Total</div>
                      <div class="text-2xl font-semibold">{{ 0 }} / {{ pointsToScore}}</div>
                    </div>
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Round</div>
                      <div class="text-2xl font-semibold">{{ 0 }}</div>
                    </div>
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Selected</div>
                      <div class="text-2xl font-semibold">{{ 0 }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Player score -->
            <div
              class="absolute bottom-4 left-4 flex flex-col gap-2 z-20"
            >
              <aside class="lg:col-span-4">
                <div class="panel">
                  <div class="font-display text-xl">Volodymyr</div>
                  <div class="mt-4 rounded-xl border border-wood-700/25 bg-parchment-50/60 p-4 space-y-3">
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Total</div>
                      <div class="text-2xl font-semibold">{{ totalScore }} / {{ pointsToScore}}</div>
                    </div>
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Round</div>
                      <div class="text-2xl font-semibold">{{ turnScore }}</div>
                    </div>
                    <div class="flex items-center justify-between">
                      <div class="text-lg text-ink-900/70">Selected</div>
                      <div class="text-2xl font-semibold">{{ lastThrowScorePreview ?? 0 }}</div>
                    </div>
                  </div>
                </div>
              </aside>
            </div>

            <!-- CANVAS CONTROLS (bottom-right) -->
            <div
              class="absolute bottom-4 right-4 flex flex-col gap-2 z-20"
            >
              <UiButton
                variant="primary"
                :disabled="!canRollNow"
                @click="roll"
              >
                {{ rolling ? 'Rolling…' : 'Roll' }}
              </UiButton>

              <UiButton
                variant="ghost"
                :disabled="!canTake"
                @click="takeSelectedAndContinue"
              >
                Take & Roll
              </UiButton>

              <UiButton
                variant="ghost"
                :disabled="!canBank"
                @click="takeSelectedAndBank"
              >
                Take & Bank
              </UiButton>
            </div>

            <!-- START OVERLAY -->
            <div
              v-if="!hasStarted"
              class="absolute inset-0 z-10 grid place-items-center bg-tavern-900/30 pointer-events-none"
            >
              <div class="text-center px-6">
                <div class="font-display text-2xl text-parchment-50/90">Press Roll to start</div>
                <div class="mt-2 text-sm text-parchment-50/60">Shake the cup and let the dice fall.</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </TavernShell>
</template>
