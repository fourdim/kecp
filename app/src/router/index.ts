import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Crate a room",
      component: () => import("@/views/CreateRoomView.vue"),
    },
    {
      path: "/join",
      name: "Join the room",
      component: () => import("@/views/JoinRoomView.vue"),
    },
    {
      path: "/create",
      redirect: "/",
    },
    {
      path: "/r/:roomID",
      name: "Room",
      component: () => import("@/views/RoomView.vue"),
    },
  ],
});

export default router;
