import {createWebHistory, createRouter, type RouteRecordRaw} from "vue-router";
import {useAuthStore} from "@/stores/auth.ts";
import {getLastCreatedGame} from "@/api/game.ts";
import {useToast} from "@/hooks/useToast.ts";

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('../pages/home/Index.vue') },
  { path: '/login', name: 'login', component: () => import('../pages/auth/Login.vue') },
  { path: '/register', name: 'register', component: () => import('../pages/auth/Registration.vue') },
  { path: '/profile', name: 'profile', component: () => import('../pages/profile/Profile.vue') },
  {
    path: '/create',
    name: 'create',
    component: () => import('../pages/lobby/CreateLobby.vue'),
    beforeEnter: async (_to, _from, next) => {
      const game = await getLastCreatedGame()

      if (game) {
        const toast = useToast()

        toast.push({ title: 'Notification', message: 'You already have a game in progress', kind: 'info' })

        return next({ name: 'lobby', params: { code: game.code } })
      }

      next()
    }
  },
  {
    path: '/join',
    name: 'join',
    component: () => import('../pages/lobby/JoinLobby.vue'),
    beforeEnter: async (_to, _from, next) => {
      const game = await getLastCreatedGame()

      if (game) return next({ name: 'lobby', params: { code: game.code } })

      next()
    }
  },
  { path: '/leaderboard', name: 'leaderboard', component: () => import('../pages/leaderboard/LeaderboardList.vue') },
  { path: '/lobby/:code', name: 'lobby', component: () => import('../pages/lobby/LobbyRoom.vue') },
  { path: '/game/:code', name: 'game', component: () => import('../pages/game/GameRoom.vue') },
  { path: '/:pathMatch(.*)*', name: 'not_found', component: () => import('@/pages/errors/NotFound.vue'),}
]

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
})

router.beforeEach((to, _from, next) => {
  const { isLoggedIn } = useAuthStore();
  const isLoginRoutes = ['login', 'register'].includes(to.name as string)

  if (isLoggedIn && isLoginRoutes) {
    next({ name: 'home' })
  } else if (!isLoggedIn && !isLoginRoutes) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router
