<script setup lang="ts">
import { onMounted, ref } from "vue";
import { RouterLink, useRouter } from "vue-router";
import Background from "@/components/BackgroundComponent.vue";
import IconLink from "@/components/icons/IconLink.vue";
import { ElMessage } from "element-plus";

const show = ref(false);
const username = ref("");
const roomID = ref("");
const router = useRouter();

onMounted(() => {
  setTimeout(() => {
    show.value = true;
  }, 200);
});

function join() {
  if (roomID.value?.length !== 16) {
    ElMessage({
      type: "info",
      message: "Invalid room ID",
    });
  }
  if (username.value === "") {
    ElMessage({
      type: "info",
      message: "Invalid username",
    });
  }
  router.push({
    name: "Room",
    params: { roomID: roomID.value, username: username.value },
  });
}
</script>

<template>
  <Background></Background>
  <div class="card-container">
    <transition name="el-fade-in-linear">
      <el-card v-show="show" :body-style="{ padding: '0px' }" class="card">
        <div class="card-title">
          <span class="title">Join the Room</span>
        </div>
        <div class="card-content">
          <el-row class="my-3">
            <el-col :xs="24" :span="6">
              <span class="h-full py-1 text-gray-600 inline-flex items-center">
                Room ID:
              </span>
            </el-col>
            <el-col :xs="24" :span="18">
              <el-input v-model="roomID" />
            </el-col>
          </el-row>
          <el-row class="my-3">
            <el-col :xs="24" :span="6">
              <span class="h-full py-1 text-gray-600 inline-flex items-center">
                Username:
              </span>
            </el-col>
            <el-col :xs="24" :span="18">
              <el-input v-model="username" />
            </el-col>
          </el-row>
          <el-row class="my-2">
            <span class="text-xs">
              Need a new room?
              <RouterLink to="/create"> <IconLink />Create one! </RouterLink>
            </span>
          </el-row>
          <el-row class="mt-8 mb-3">
            <el-button type="primary" class="mx-auto" @click="join">
              <span class="px-5 text-base">Join</span>
            </el-button>
          </el-row>
        </div>
      </el-card>
    </transition>
  </div>
</template>

<style>
.card-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.card-title {
  display: flex;
  justify-content: center;
  padding: 30px;
}

.title {
  font-size: x-large;
  font-weight: 600;
}

.card-content {
  background-color: #ecf5ff;
  padding: 30px 50px;
}

.card {
  width: 480px;
}

@media (max-width: 768px) {
  .card-content {
    padding: 20px;
  }

  .card {
    width: 70vw;
  }
}

a {
  color: #0067b8;
}
</style>
