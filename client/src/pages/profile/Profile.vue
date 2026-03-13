<script setup lang="ts">
import { computed, ref } from "vue";
import TavernShell from "@/components/wrappers/TavernShell.vue";
import UiButton from "@/components/form/UiButton.vue";
import { useAuthStore } from "@/stores/auth.ts";

const auth = useAuthStore();
const user = computed(() => auth.user);

// demo data (потім з бекенду)
const gold = ref(1000);
const rank = ref("Tavern Regular");
const title = ref("Dice Drifter");
const joined = ref("Feb 2026");
const lastSeen = ref("Today");
const favorite = ref("Oak dice set");
const bestCombo = ref("Straight (1–5)");

const winRate = ref(52);
const gamesPlayed = ref(128);
const gamesWon = ref(67);

const initials = computed(() => {
  const u = user.value?.username ?? "Player";
  return u.slice(0, 2).toUpperCase();
});

const kpi = computed(() => [
  { label: "Gold", value: `${gold.value}`, hint: "Balance" },
  { label: "Games", value: `${gamesPlayed.value}`, hint: "Total played" },
  { label: "Wins", value: `${gamesWon.value}`, hint: "All-time" },
  { label: "Win rate", value: `${winRate.value}%`, hint: "Overall" },
]);

const achievements = computed(() => [
  { title: "Lucky Hand", desc: "Score 1,000+ points in a single round", done: true },
  { title: "Cold Streak", desc: "Lose 3 rounds in a row", done: true },
  { title: "High Roller", desc: "Hold 500+ gold at once", done: false },
  { title: "Tavern Legend", desc: "Win 50 matches", done: gamesWon.value >= 50 },
]);

const recentMatches = computed(() => [
  { result: "Win", score: "820", vs: "Miller_John", time: "2h ago" },
  { result: "Loss", score: "410", vs: "SirRadzig", time: "Yesterday" },
  { result: "Win", score: "600", vs: "Theresa", time: "3d ago" },
]);
</script>

<template>
  <TavernShell>
    <div class="aaa mx-auto w-full max-w-6xl">
      <!-- HERO / CHARACTER CARD -->
      <section class="hero">
        <div class="hero-inner">
          <div class="hero-left">
            <div class="portrait">
              <div class="portrait-ring"></div>
              <div class="portrait-core">
                <span class="font-display text-xl tracking-wide">{{ initials }}</span>
              </div>
            </div>

            <div class="hero-meta">
              <div class="name-row">
                <div class="nameplate">
                  <div class="nameplate-name font-display">
                    {{ user?.username ?? "Player" }}
                  </div>
                </div>

                <span class="rank-chip">{{ rank }}</span>
              </div>

              <div class="subline">
                <span class="sub-k">Title:</span>
                <span class="sub-v">{{ title }}</span>
                <span class="dot"></span>
                <span class="sub-k">Last seen:</span>
                <span class="sub-v">{{ lastSeen }}</span>
              </div>

              <div class="actions">
                <UiButton variant="primary" @click="$router.push('/')">Back to tavern</UiButton>
                <UiButton variant="ghost" @click="auth.logout()">Logout</UiButton>
              </div>
            </div>
          </div>

          <div class="hero-right">
            <div class="gold-box">
              <div class="gold-k">Balance</div>
              <div class="gold-v tabular-nums">
                {{ gold }}
                <span class="coin" title="Gold"></span>
              </div>
              <div class="gold-sub">Spend wisely, gambler.</div>
            </div>

            <div class="mini-grid">
              <div class="mini">
                <div class="mini-k">Joined</div>
                <div class="mini-v">{{ joined }}</div>
              </div>
              <div class="mini">
                <div class="mini-k">Favorite</div>
                <div class="mini-v">{{ favorite }}</div>
              </div>
              <div class="mini">
                <div class="mini-k">Best combo</div>
                <div class="mini-v">{{ bestCombo }}</div>
              </div>
              <div class="mini">
                <div class="mini-k">Unlocked</div>
                <div class="mini-v">{{ achievements.filter(a => a.done).length }}/{{ achievements.length }}</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- GRID -->
      <div class="grid mt-6 grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- STATS -->
        <section class="panel lg:col-span-1">
          <div class="panel-head">
            <div class="orn-title">Quick stats</div>
            <div class="panel-sub">All-time</div>
          </div>

          <div class="kpi-grid">
            <div v-for="item in kpi" :key="item.label" class="kpi">
              <div class="kpi-top">
                <div class="kpi-label">{{ item.label }}</div>
                <div class="kpi-hint">{{ item.hint }}</div>
              </div>

              <div class="kpi-value tabular-nums">
                <span v-if="item.label === 'Gold'" class="inline-flex items-center gap-2">
                  {{ item.value }}
                  <span class="coin coin-sm"></span>
                </span>
                <span v-else>{{ item.value }}</span>
              </div>

              <div class="kpi-bar">
                <div class="kpi-fill"></div>
              </div>
            </div>
          </div>

          <div class="footnote">
            Next title: <span class="accent">“Oak Ring”</span>
            <div class="xp">
              <div class="xp-fill" style="width:62%"></div>
            </div>
            <div class="xp-sub">62% to level 2</div>
          </div>
        </section>

        <!-- ACHIEVEMENTS -->
        <section class="panel lg:col-span-2">
          <div class="panel-head">
            <div class="orn-title">Achievements</div>
            <div class="panel-sub">{{ achievements.filter(a => a.done).length }}/{{ achievements.length }}</div>
          </div>

          <div class="ach-grid">
            <div
              v-for="a in achievements"
              :key="a.title"
              class="ach"
              :class="a.done ? 'ach-done' : 'ach-locked'"
            >
              <div class="ach-badge" :class="a.done ? 'ach-badge-done' : 'ach-badge-locked'">
                <span v-if="a.done">✓</span>
                <span v-else>🔒</span>
              </div>

              <div class="ach-body">
                <div class="ach-title">{{ a.title }}</div>
                <div class="ach-desc">{{ a.desc }}</div>
              </div>

              <div class="ach-glow"></div>
            </div>
          </div>

          <div class="tip">
            Tip: later we can reward frames, dice skins, and tavern titles for these.
          </div>
        </section>

        <!-- RECENT MATCHES -->
        <section class="panel lg:col-span-3">
          <div class="panel-head">
            <div class="orn-title">Recent matches</div>
            <UiButton variant="ghost" class="!px-3 !py-2 text-sm" @click="$router.push('/')">
              Play
            </UiButton>
          </div>

          <div class="matches">
            <div v-for="m in recentMatches" :key="m.time + m.vs" class="match">
              <div class="match-left">
                <span class="pill" :class="m.result === 'Win' ? 'pill-win' : 'pill-loss'">
                  {{ m.result }}
                </span>
                <div class="match-vs">vs <span class="match-name">{{ m.vs }}</span></div>
              </div>

              <div class="match-right">
                <div class="match-score tabular-nums">{{ m.score }}</div>
                <div class="match-time">{{ m.time }}</div>
              </div>

              <div class="match-line"></div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </TavernShell>
</template>

<style scoped>
/* ===== AAA vibe wrapper ===== */
.aaa {
  --gold: rgba(212,145,57,1);
  --goldSoft: rgba(212,145,57,.24);
  --ink: rgba(12,9,6,1);
  --parch: rgba(255,240,210,1);
}

/* ===== HERO CARD ===== */
.hero {
  border-radius: 1.4rem;
  border: 1px solid rgba(255,240,210,.14);
  background:
    radial-gradient(circle at 18% 18%, rgba(212,145,57,.20), transparent 55%),
    radial-gradient(circle at 90% 40%, rgba(255,240,210,.08), transparent 45%),
    linear-gradient(180deg, rgba(18,14,9,.92), rgba(10,8,6,.92));
  box-shadow: 0 20px 60px rgba(0,0,0,.45);
  overflow: hidden;
  position: relative;
}

/* subtle "noise" using layered gradients (no images) */
.hero::before{
  content:"";
  position:absolute;
  inset:0;
  background:
    radial-gradient(circle at 20% 20%, rgba(255,240,210,.08), transparent 55%),
    radial-gradient(circle at 70% 60%, rgba(212,145,57,.08), transparent 60%);
  opacity:.28;
  pointer-events:none;
}

.hero::after{
  content:"";
  position:absolute;
  inset:-2px;
  background: radial-gradient(circle at 50% 0%, rgba(0,0,0,0), rgba(0,0,0,.55) 70%);
  pointer-events:none;
}

.hero-inner{
  position: relative;
  display: grid;
  grid-template-columns: 1.3fr 0.9fr;
  gap: 1.25rem;
  padding: 1.25rem;
}

.hero-inner::before{
  content:"";
  position:absolute;
  inset:0;
  background:
    linear-gradient(90deg, rgba(0,0,0,.40) 0%, rgba(0,0,0,.18) 42%, transparent 70%);
  pointer-events:none;
}

@media (max-width: 1024px){
  .hero-inner{ grid-template-columns: 1fr; }
}

.hero-left{
  display:flex;
  gap: 1rem;
  align-items: center;
}
.hero-meta{ min-width: 0; flex: 1; }

.portrait{
  width: 74px;
  height: 74px;
  position: relative;
  flex: 0 0 auto;
}
.portrait-ring{
  position:absolute; inset:-2px;
  border-radius: 1.35rem;
  background:
    radial-gradient(circle at 30% 30%, rgba(255,240,210,.22), transparent 55%),
    linear-gradient(135deg, rgba(212,145,57,.35), rgba(255,240,210,.08));
  filter: blur(.2px);
  box-shadow:
    0 16px 28px rgba(0,0,0,.35),
    inset 0 1px 0 rgba(255,240,210,.12);
}
.portrait-core{
  position:absolute; inset:6px;
  border-radius: 1.15rem;
  display:grid; place-items:center;
  border: 1px solid rgba(255,240,210,.18);
  background:
    radial-gradient(circle at 30% 30%, rgba(212,145,57,.20), transparent 60%),
    rgba(255,240,210,.06);
  color: rgba(255,240,210,.92);
}

.name-row{
  display:flex;
  flex-wrap: wrap;
  align-items: center;
  gap: .6rem;
}

.nameplate{
  position: relative;
  display:inline-flex;
  padding: .45rem .85rem;
  border-radius: 999px;
  border: 1px solid rgba(255,240,210,.14);
  background: linear-gradient(180deg, rgba(255,240,210,.07), rgba(255,240,210,.03));
  box-shadow:
    inset 0 1px 0 rgba(255,240,210,.08),
    0 10px 24px rgba(0,0,0,.35);
}
.nameplate::before{
  content:"";
  position:absolute;
  inset:-10px -18px;
  border-radius: 999px;
  background: radial-gradient(circle at 30% 50%, rgba(212,145,57,.30), transparent 60%);
  filter: blur(12px);
  opacity: .95;
  pointer-events:none;
}
.nameplate-name{
  position: relative;
  font-size: 1.35rem;
  letter-spacing: .06em;
  color: rgba(255,232,196,.96);
  text-shadow:
    0 1px 0 rgba(0,0,0,.55),
    0 0 14px rgba(212,145,57,.20);
}

.rank-chip{
  font-size: .75rem;
  padding: .35rem .6rem;
  border-radius: 999px;
  border: 1px solid rgba(255,240,210,.12);
  background: rgba(255,240,210,.05);
  color: rgba(255,240,210,.78);
}
.subline{ color: rgba(255,240,210,.72); }
.sub-v{ color: rgba(255,240,210,.92); }
.rank-chip{ color: rgba(255,240,210,.86); }

.sub-k{ color: rgba(255,240,210,.50); text-transform: uppercase; letter-spacing: .08em; font-size: .7rem; }
.dot{
  width: 4px; height: 4px;
  border-radius: 999px;
  background: rgba(212,145,57,.55);
  opacity: .8;
}

.actions{
  margin-top: .9rem;
  display:flex;
  flex-wrap: wrap;
  gap: .5rem;
}

.hero-right{
  display:flex;
  flex-direction: column;
  gap: .8rem;
}

.gold-box{
  border-radius: 1.1rem;
  border: 1px solid rgba(255,240,210,.12);
  background:
    radial-gradient(circle at 20% 20%, rgba(212,145,57,.18), transparent 55%),
    rgba(255,240,210,.03);
  padding: .9rem 1rem;
  box-shadow: inset 0 1px 0 rgba(255,240,210,.06);
}
.gold-k{
  font-size: .7rem;
  text-transform: uppercase;
  letter-spacing: .12em;
  color: rgba(255,240,210,.55);
}
.gold-v{
  margin-top: .25rem;
  font-size: 1.55rem;
  font-weight: 700;
  color: rgba(255,240,210,.92);
  display:flex;
  align-items:center;
  gap: .5rem;
  text-shadow: 0 0 18px rgba(212,145,57,.18);
}
.gold-sub{
  margin-top: .25rem;
  font-size: .8rem;
  color: rgba(255,240,210,.55);
}

.mini-grid{
  display:grid;
  grid-template-columns: 1fr 1fr;
  gap: .6rem;
}
.mini{
  border-radius: .9rem;               /* менше “pill”, більше AAA */
  border: 1px solid rgba(255,240,210,.10);
  background: rgba(255,240,210,.03);
  padding: .7rem .8rem;
}
.mini-k{
  font-size: .65rem;
  text-transform: uppercase;
  letter-spacing: .12em;
  color: rgba(255,240,210,.52);
}
.mini-v{
  margin-top: .2rem;
  color: rgba(255,240,210,.90);
}

/* ===== PANELS ===== */
.panel{
  border-radius: 1.25rem;
  border: 1px solid rgba(255,240,210,.12);
  background:
    radial-gradient(circle at 25% 0%, rgba(212,145,57,.12), transparent 55%),
    linear-gradient(180deg, rgba(18,14,9,.90), rgba(10,8,6,.90));
  box-shadow: 0 16px 40px rgba(0,0,0,.38);
  padding: 1rem 1.1rem;
  position: relative;
  overflow: hidden;
}
.panel::before{
  content:"";
  position:absolute; inset:0;
  background:
    repeating-linear-gradient(135deg, rgba(255,240,210,.02) 0 2px, transparent 2px 7px);
  opacity:.25;
  pointer-events:none;
}

.panel-head{
  display:flex;
  align-items:center;
  justify-content: space-between;
  gap: .75rem;
  position: relative;
}
.panel-sub{
  font-size: .75rem;
  color: rgba(255,240,210,.55);
}

/* Ornaments */
.orn-title{
  position: relative;
  font-size: 1.05rem;
  font-weight: 800;
  letter-spacing: .14em;
  text-transform: uppercase;
  color: rgba(255,240,210,.92);
  padding-left: 1.15rem;
  text-shadow: 0 1px 0 rgba(0,0,0,.55), 0 0 14px rgba(212,145,57,.10);
}

.orn-title::before{
  width: 12px;
  height: 12px;
  border-color: rgba(212,145,57,.70);
}

.actions :deep(button){
  box-shadow: 0 14px 28px rgba(0,0,0,.28);
}

/* ===== KPI ===== */
.kpi-grid{
  margin-top: 1rem;
  display:grid;
  grid-template-columns: 1fr 1fr;
  gap: .75rem;
}
.kpi{
  border-radius: .9rem;
  border: 1px solid rgba(255,240,210,.10);
  background: rgba(255,240,210,.03);
  padding: .8rem .85rem;
  position: relative;
  transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
}
.kpi:hover{
  transform: translateY(-1px);
  border-color: rgba(212,145,57,.22);
  box-shadow: 0 18px 30px rgba(0,0,0,.22), inset 0 0 0 1px rgba(212,145,57,.12);
}
.kpi-top{
  display:flex;
  justify-content: space-between;
  align-items: baseline;
  gap: .5rem;
}
.kpi-label{
  font-size: .7rem;
  text-transform: uppercase;
  letter-spacing: .12em;
  color: rgba(255,240,210,.55);
}
.kpi-hint{
  font-size: .7rem;
  color: rgba(255,240,210,.42);
}
.kpi-value{
  margin-top: .35rem;
  font-size: 1.25rem;
  font-weight: 700;
  color: rgba(255,240,210,.92);
  text-shadow: 0 0 14px rgba(212,145,57,.12);
}
.kpi-bar{
  margin-top: .55rem;
  height: 6px;
  border-radius: 999px;
  background: rgba(255,240,210,.08);
  overflow:hidden;
}
.kpi-fill{
  height: 100%;
  width: 62%;
  background: linear-gradient(90deg, rgba(212,145,57,.25), rgba(212,145,57,.55));
}

.footnote{
  margin-top: 1rem;
  font-size: .85rem;
  color: rgba(255,240,210,.70);
}
.accent{ color: rgba(212,145,57,.95); }
.xp{
  margin-top: .5rem;
  height: 8px;
  border-radius: 999px;
  background: rgba(255,240,210,.08);
  overflow:hidden;
}
.xp-fill{
  height:100%;
  background: linear-gradient(90deg, rgba(212,145,57,.25), rgba(212,145,57,.60));
}
.xp-sub{
  margin-top: .35rem;
  font-size: .75rem;
  color: rgba(255,240,210,.50);
}

/* ===== Achievements ===== */
.ach-grid{
  margin-top: 1rem;
  display:grid;
  grid-template-columns: 1fr 1fr;
  gap: .75rem;
}
@media (max-width: 640px){
  .ach-grid{ grid-template-columns: 1fr; }
}
.ach{
  border-radius: 1rem;
  border: 1px solid rgba(255,240,210,.10);
  background: rgba(255,240,210,.03);
  padding: .9rem;
  display:flex;
  gap: .8rem;
  align-items: flex-start;
  position: relative;
  overflow:hidden;
  transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
}
.ach:hover{
  transform: translateY(-1px);
  box-shadow: 0 18px 30px rgba(0,0,0,.22);
}
.ach-badge{
  width: 36px; height: 36px;
  border-radius: .9rem;
  display:grid;
  place-items:center;
  border: 1px solid rgba(255,240,210,.12);
  flex: 0 0 auto;
}
.ach-badge-done{
  background: rgba(212,145,57,.18);
  color: rgba(255,240,210,.92);
  border-color: rgba(212,145,57,.28);
}
.ach-badge-locked{
  background: rgba(255,240,210,.05);
  color: rgba(255,240,210,.62);
}
.ach-body{ min-width:0; }
.ach-title{
  color: rgba(255,240,210,.94);
  font-weight: 700;
}
.ach-desc{
  margin-top: .2rem;
  color: rgba(255,240,210,.58);
  font-size: .9rem;
}
.ach-done{
  border-color: rgba(212,145,57,.22);
}
.ach-locked{
  opacity: .88;
}
.ach-glow{
  position:absolute;
  inset:-30px -60px auto -60px;
  height: 120px;
  background: radial-gradient(circle at 30% 50%, rgba(212,145,57,.18), transparent 60%);
  filter: blur(14px);
  opacity: .0;
  transition: opacity 180ms ease;
  pointer-events:none;
}
.ach:hover .ach-glow{ opacity: .9; }

.tip{
  margin-top: .9rem;
  font-size: .8rem;
  color: rgba(255,240,210,.55);
}

/* ===== Matches ===== */
.matches{
  margin-top: 1rem;
  display:grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: .75rem;
}
@media (max-width: 1024px){
  .matches{ grid-template-columns: 1fr; }
}
.match{
  border-radius: 1rem;
  border: 1px solid rgba(255,240,210,.10);
  background: rgba(255,240,210,.03);
  padding: .9rem;
  display:flex;
  justify-content: space-between;
  gap: 1rem;
  position: relative;
  overflow:hidden;
  transition: transform 140ms ease, box-shadow 140ms ease, border-color 140ms ease;
}
.match:hover{
  transform: translateY(-1px);
  border-color: rgba(212,145,57,.18);
  box-shadow: 0 18px 30px rgba(0,0,0,.22);
}
.match-left{
  display:flex;
  align-items:center;
  gap: .65rem;
  min-width:0;
}
.match-vs{
  color: rgba(255,240,210,.78);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.match-name{ color: rgba(255,240,210,.94); font-weight: 700; }
.match-right{
  text-align:right;
  flex: 0 0 auto;
}
.match-score{
  color: rgba(255,240,210,.92);
  font-weight: 800;
}
.match-time{
  margin-top: .2rem;
  color: rgba(255,240,210,.52);
  font-size: .8rem;
}
.match-line{
  position:absolute;
  inset:auto 0 0 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(212,145,57,.45), transparent);
  opacity: .35;
}

/* Pills */
.pill{
  font-size: .72rem;
  padding: .28rem .55rem;
  border-radius: 999px;
  border: 1px solid rgba(255,240,210,.12);
  letter-spacing: .06em;
  text-transform: uppercase;
}
.pill-win{
  background: rgba(76,175,80,.14);
  color: rgba(210,255,215,.9);
}
.pill-loss{
  background: rgba(244,67,54,.12);
  color: rgba(255,210,210,.9);
}

/* Coin */
.coin{
  width: 18px;
  height: 18px;
  display:inline-block;
  border-radius: 999px;
  border: 1px solid rgba(212,145,57,.55);
  background:
    radial-gradient(circle at 35% 35%, rgba(255,240,210,.45), transparent 55%),
    radial-gradient(circle at 65% 70%, rgba(212,145,57,.45), transparent 60%),
    rgba(212,145,57,.18);
  box-shadow: 0 8px 18px rgba(212,145,57,.16);
}
.coin-sm{ width: 16px; height: 16px; }

.tabular-nums{ font-variant-numeric: tabular-nums; }
</style>
