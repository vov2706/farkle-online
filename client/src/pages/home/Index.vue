<script setup lang="ts">
import TavernShell from '../../components/TavernShell.vue'
import MenuLink from '../../components/MenuLink.vue'
import {useAuthStore} from "@/stores/auth.ts";

const auth = useAuthStore();
const currentGame = auth.user?.current_game
</script>

<template>
  <TavernShell>
    <div class="mx-auto max-w-2xl space-y-3">
      <template v-if="!currentGame">
        <MenuLink
          title="Create game"
          description="Host a room and share the code with friends."
          to="/create"
        />

        <MenuLink
          title="Join room"
          description="Enter a room code or pick from public rooms."
          to="/join"
        />
      </template>

      <MenuLink
        v-else
        title="Return to the lobby"
        description="Return to the joined game."
        :to="`/lobby/${currentGame.code}`"
      />

<!--      <MenuLink-->
<!--        title="Leaderboard"-->
<!--        description="Top players, wins and rating."-->
<!--        to="/leaderboard"-->
<!--      />-->
    </div>
  </TavernShell>
</template>
