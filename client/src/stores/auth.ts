import { defineStore } from "pinia";
import fetchApi from "@/packages/fetchApi.ts";
import {getProfile, type User} from "@/api/user.ts";
import router from "@/routes";
import {useToast} from "@/hooks/useToast.ts";

const ACCESS_TOKEN = "access_token";

interface LoginResponse {
  message: string
  token: string
}

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: (null as User | null),
    token: localStorage.getItem(ACCESS_TOKEN) || "",
  }),

  getters: {
    isLoggedIn: (s) => !!s.token,
  },

  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem(ACCESS_TOKEN, token);
    },

    setUser(user: User | null) {
      this.user = user;
    },

    logout() {
      this.token = "";
      this.user = null;
      localStorage.removeItem(ACCESS_TOKEN);

      useToast().push({title: "Success", message: "Logged out", kind: "success"});

      (async () => await router.push("/login"))();
    },

    async login(username: string, password: string) {
      const {data} = await fetchApi.post<LoginResponse>("/login", {username, password});
      this.setToken(data.token);
      await this.fetchProfile();
      await router.push("/");
    },

    async register(username: string, password: string) {
      const {data} = await fetchApi.post<LoginResponse>("/register", {username, password});
      this.setToken(data.token);
      await this.fetchProfile();
      await router.push("/");
    },

    async fetchProfile() {
      const user = await getProfile();
      this.setUser(user);
      return user;
    },
    async bootstrap() {
      try {
        if (!this.token) return;
        await this.fetchProfile();
      } catch (e) {
        this.logout();
      }
    },
  },
});
