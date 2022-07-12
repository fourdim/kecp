<script setup lang="ts">
import logo from "@/assets/logo.svg";
import IconSetting from "./icons/IconSetting.vue";
</script>

<template>
  <div class="flex justify-between bg-white py-4 px-5">
    <RouterLink to="/create">
      <div class="flex">
        <img :src="logo" class="w-10" />
        <span class="text-4xl text-black <lg:hidden font-bold mx-5">Kecp</span>
      </div>
    </RouterLink>
    <div class="mr-10" @click="$emit('open-settings')">
      <IconSetting
        class="p-1 h-full rounded-1/2 hover:shadow-inner"
      ></IconSetting>
    </div>
  </div>
</template>
