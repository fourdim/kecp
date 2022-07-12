<script setup lang="ts">
import { onMounted, ref } from "vue";
import videojs, { type VideoJsPlayer } from "video.js";
import { RouterLink, useRoute, useRouter } from "vue-router";
import {
  KecpSignal,
  KecpRoom,
  AnswerPeer,
  KecpEventType,
  KecpMessageType,
  type KecpMessage,
} from "kecp-webrtc";
import {
  ElMessage,
  ElMessageBox,
  ElScrollbar,
  type UploadFile,
  type UploadFiles,
} from "element-plus";
import { kecpEndpoint } from "@/config";
import Background from "@/components/BackgroundComponent.vue";
import Header from "@/components/HeaderComponent.vue";
import IconLink from "@/components/icons/IconLink.vue";

type KecpChat = KecpMessage & { chatID: number };
interface HTMLMediaElementWithCaputreStream extends HTMLMediaElement {
  captureStream(): MediaStream;
}

const route = useRoute();
const router = useRouter();
const videoPlayerEl = ref<HTMLMediaElementWithCaputreStream>();
const player = ref<VideoJsPlayer>();
const roomID = ref("");
const room = ref<KecpRoom>();
const username = ref("");
const usernameInput = ref("");
const userList = ref<string[]>([]);
const targetUser = ref<string>("");
const chatList = ref<KecpChat[]>([]);
const sig = new KecpSignal(kecpEndpoint);
const connected = ref(false);
const sendValue = ref("");
const chatIDCouter = ref(0);
const chatBox = ref<HTMLElement[]>([]);
const fileURL = ref("");
const scrollBar = ref<InstanceType<typeof ElScrollbar>>();
const innerRef = ref<HTMLDivElement>();
const settingsDialogVisible = ref(false);
const initDialogVisible = ref(false);
const initDialogLoading = ref(false);

function connect() {
  if (room.value === undefined) {
    return;
  }

  room.value.connect(username.value);

  room.value.on(KecpEventType.Open, () => {
    connected.value = true;
  });

  room.value.on(KecpEventType.UserListInit, (evt: CustomEvent<string[]>) => {
    userList.value = evt.detail.slice();
    initDialogLoading.value = false;
    initDialogVisible.value = false;
  });

  room.value.on(KecpEventType.UserJoin, (evt: CustomEvent<string>) => {
    userList.value.push(evt.detail);
  });

  room.value.on(KecpEventType.UserLeave, (evt: CustomEvent<string>) => {
    if (userList.value.indexOf(evt.detail) !== -1) {
      userList.value.splice(userList.value.indexOf(evt.detail));
    }
  });

  room.value.on(KecpEventType.Chat, (evt: CustomEvent<KecpMessage>) => {
    chatIDCouter.value += 1;
    chatList.value.push({
      chatID: chatIDCouter.value,
      type: evt.detail.type,
      name: evt.detail.name,
      target: evt.detail.target,
      payload: evt.detail.payload,
    });
    setTimeout(() => {
      scrollBar.value?.scrollTo(0, innerRef.value?.clientHeight);
    }, 150);
  });

  room.value.on(
    KecpEventType.VideoOffer,
    async (evt: CustomEvent<AnswerPeer>) => {
      const stream = new MediaStream();
      evt.detail.setHandleTrackEvent((event) => {
        stream.addTrack(event.track);
        if (videoPlayerEl.value) {
          if (videoPlayerEl.value.srcObject !== stream) {
            videoPlayerEl.value.srcObject = stream;
            videoPlayerEl.value.autoplay = true;
          }
        }
      });
      await evt.detail.preAnswer();
      ElMessageBox.confirm(
        `${evt.detail.getTargetName()} will establish a connection with you. Confirm?`,
        "Warning",
        {
          confirmButtonText: "Confirm",
          cancelButtonText: "Cancel",
          type: "warning",
        }
      )
        .then(async () => {
          await evt.detail.preAnswerConfirm();
          ElMessage({
            type: "success",
            message: "Offer completed",
          });
        })
        .catch(() => {
          evt.detail.close();
          if (videoPlayerEl.value) {
            videoPlayerEl.value.srcObject = null;
            videoPlayerEl.value.autoplay = false;
          }
          ElMessage({
            type: "info",
            message: "Offer canceled",
          });
        });
    }
  );

  room.value.on(KecpEventType.VideoAnswer, (evt: CustomEvent<KecpMessage>) => {
    if (evt.detail.payload.type === "answer") {
      ElMessage({
        type: "success",
        message: "Offer completed",
      });
    }
  });

  room.value.on(KecpEventType.Error, (evt: CustomEvent<KecpMessage>) => {
    if (evt.detail.payload === "cannot join the room") {
      router.replace("/create");
    }
    switch (evt.detail.payload) {
      case "name is already in use":
        ElMessage({
          type: "error",
          message: "Name is already in use. Try another one.",
        });
        username.value = "";
        connected.value = false;
        initDialogLoading.value = false;
        break;
      case "not a valid name":
        ElMessage({
          type: "error",
          message: "Not a valid name. Try another one.",
        });
        username.value = "";
        connected.value = false;
        initDialogLoading.value = false;
        break;
    }
  });
}

onMounted(() => {
  roomID.value = route.params["roomID"] as string;
  room.value = sig.getRoom({
    roomID: roomID.value,
  });
  username.value = (route.params["username"] ?? "") as string;
  if (roomID.value.length !== 16) {
    router.replace("/create");
  }
  if (videoPlayerEl.value !== undefined) {
    player.value = videojs(videoPlayerEl.value, {
      controls: true,
      fill: true,
    });
  }
  if (!connected.value && username.value !== "") {
    connect();
  } else if (username.value === "") {
    initDialogVisible.value = true;
  }
});

function uploadFile(uploadFile: UploadFile, uploadFiles: UploadFiles) {
  if (uploadFiles === null) {
    return;
  }
  const file = uploadFile.raw;
  if (file !== undefined) {
    if (fileURL.value !== "") {
      URL.revokeObjectURL(fileURL.value);
    }
    fileURL.value = URL.createObjectURL(file);
    if (videoPlayerEl.value !== null && videoPlayerEl.value !== undefined) {
      videoPlayerEl.value.srcObject = null;
      player.value?.src({ type: file.type, src: fileURL.value });
    }
  }
}

function removeFile() {
  if (videoPlayerEl.value !== null && videoPlayerEl.value !== undefined) {
    videoPlayerEl.value.srcObject = null;
    player.value?.reset();
  }
  if (fileURL.value !== "") {
    URL.revokeObjectURL(fileURL.value);
  }
  fileURL.value = "";
}

function exceedFile() {
  ElMessage({
    type: "warning",
    message: "Limit to 1 file, please remove the existing one",
  });
}

function handleSendButton() {
  if (sendValue.value === "") {
    return;
  }
  room.value?.send(
    JSON.stringify({
      type: KecpMessageType.Chat,
      name: room.value.getSelfName(),
      payload: sendValue.value,
    })
  );
}

function handleJoinButton() {
  username.value = usernameInput.value;
  initDialogLoading.value = true;
  connect();
}

function handleDialButton() {
  if (targetUser.value === room.value?.getSelfName()) {
    return;
  }

  const stream = videoPlayerEl.value?.captureStream();
  if (stream) {
    const offer = room.value?.newOffer(targetUser.value);
    offer?.addStream(stream);
    offer?.setBandwidth(10000);
  }
}

function playerRewind() {
  if (player.value !== undefined) {
    const target = player.value.currentTime() - 5;
    player.value.currentTime(target);
  }
}

function playerForward() {
  if (player.value !== undefined) {
    const target = player.value.currentTime() + 5;
    player.value.currentTime(target);
  }
}

function playerPause() {
  if (player.value !== undefined) {
    if (player.value.paused()) {
      player.value.play();
    } else {
      player.value.pause();
    }
  }
}
</script>

<template>
  <Background></Background>
  <Header @open-settings="settingsDialogVisible = true"></Header>
  <div
    class="room-grid-container"
    @keydown.arrow-left="playerRewind"
    @keydown.arrow-right="playerForward"
    @keydown.space="playerPause"
  >
    <div class="left flex flex-col">
      <div class="video-container bg-white">
        <video
          ref="videoPlayerEl"
          class="video-js vjs-big-play-centered vjs-show-big-play-button-on-pause"
        ></video>
      </div>
      <div
        class="chat-input-container bg-white"
        @keyup.enter="handleSendButton()"
      >
        <span>Chat:</span>
        <el-input class="chat-input px-2" v-model="sendValue" />
        <el-button id="send" :disabled="!connected" @click="handleSendButton()">
          Send
        </el-button>
      </div>
    </div>
    <div class="right">
      <el-scrollbar class="chat-container bg-white" ref="scrollBar">
        <div ref="innerRef">
          <span
            class="leading-snug"
            v-for="chat in chatList"
            ref="chatBox"
            :key="chat.chatID"
          >
            {{ chat.name ?? "" }}: {{ chat.payload }} <br />
          </span>
        </div>
      </el-scrollbar>
    </div>
  </div>
  <el-dialog v-model="settingsDialogVisible" custom-class="setting-dialog">
    <div class="card-title">
      <span class="title">Settings</span>
    </div>
    <div class="ctrl-container flex-shrink-0 px-12 py-8">
      <span class="text-lg mb-2">Choose a video file:</span>
      <div class="mb-2">
        <el-upload
          drag
          :auto-upload="false"
          :on-change="uploadFile"
          :on-remove="removeFile"
          :on-exceed="exceedFile"
          :limit="1"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            Drop file here or <em>click to upload</em>
          </div>
        </el-upload>
      </div>
      <span class="text-lg mb-2">Select a target to dial:</span>
      <div class="flex">
        <el-select v-model="targetUser" placeholder="Select">
          <el-option
            v-for="item in userList"
            :key="item"
            :label="item"
            :value="item"
          />
        </el-select>
        <el-button
          id="send"
          class="mx-3"
          :disabled="!connected || targetUser === '' || fileURL == ''"
          @click="handleDialButton()"
        >
          Dial
        </el-button>
      </div>
    </div>
  </el-dialog>
  <el-dialog
    :show-close="false"
    :close-on-press-escape="false"
    v-model="initDialogVisible"
    custom-class="init-dialog"
    top="30vh"
  >
    <div class="card-title">
      <span class="title">Join the Room</span>
    </div>
    <div class="card-content">
      <el-row class="my-3">
        <el-col :xs="24" :span="6">
          <span class="h-full py-1 text-gray-600 inline-flex items-center">
            Username:
          </span>
        </el-col>
        <el-col :xs="24" :span="18">
          <el-input v-model="usernameInput" />
        </el-col>
      </el-row>
      <el-row class="my-2">
        <span class="text-xs">
          Need a new room?
          <RouterLink to="/create"> <IconLink />Create one! </RouterLink>
        </span>
      </el-row>
      <el-row class="mt-8 mb-3">
        <el-button type="primary" class="mx-auto" @click="handleJoinButton">
          <span class="px-5 text-base">Join</span>
        </el-button>
      </el-row>
    </div>
  </el-dialog>
</template>

<style>
@import "video.js/dist/video-js.css";
</style>

<style scoped>
.room-grid-container {
  padding: 30px 70px;
  display: grid;
  height: 100%;
  grid-template-areas: "left right";
  grid-template-columns: auto 300px;
  grid-auto-rows: 1fr;
  gap: 0 50px;
}

@media (max-width: 1024px) {
  .room-grid-container {
    padding: 30px 4vw;
    grid-template-areas:
      "left"
      "right";
    grid-template-columns: auto;
    grid-template-rows: auto 40vh;
  }
}

.left {
  grid-area: left;
}

.video-container {
  aspect-ratio: 16/9;
  box-shadow: var(--el-box-shadow-light);
}

.chat-input-container {
  box-shadow: var(--el-box-shadow-light);
  border-bottom-left-radius: var(--el-border-radius-base);
  border-bottom-right-radius: var(--el-border-radius-base);
  display: flex;
  padding: 0 20px;
  align-items: center;
  height: 50px;
}

.chat-input {
  width: 100%;
}

.ctrl-container {
  grid-area: ctrl-container;
  display: flex;
  flex-direction: column;
  background-color: #ecf5ff;
}

.right {
  grid-area: right;
  display: flex;
  flex-direction: column;
}

@media (max-width: 1024px) {
  .right {
    margin-top: 30px;
  }
}

.chat-container {
  padding: 16px;
  border-radius: var(--el-border-radius-base);
  box-shadow: var(--el-box-shadow-light);
  max-height: 70vh;
}

.card-title {
  display: flex;
  justify-content: center;
  padding: 25px;
  padding-top: 10px;
}

.title {
  font-size: x-large;
  font-weight: 600;
}

.card-content {
  background-color: #ecf5ff;
  padding: 30px 50px;
}

a {
  color: #0067b8;
}
</style>

<style>
.el-dialog__body {
  padding: 0px;
}

@media (max-width: 1024px) {
  .setting-dialog {
    --el-dialog-width: 70%;
  }
}

.init-dialog {
  --el-dialog-width: 30%;
}

@media (max-width: 1536px) {
  .init-dialog {
    --el-dialog-width: 40%;
  }
}

@media (max-width: 1024px) {
  .init-dialog {
    --el-dialog-width: 60%;
  }
}

@media (max-width: 768px) {
  .init-dialog {
    --el-dialog-width: 70%;
  }
}
</style>
