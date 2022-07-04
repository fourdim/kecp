<script setup lang="ts">
// eslint-disable-next-line import/no-extraneous-dependencies
import { ref } from 'vue';
import {
  KecpSignal, KecpConnection, KecpEventType, KecpMessageType, KecpMessage, AnswerPeer,
} from '../src/main';

interface HTMLMediaElementWithCaputreStream extends HTMLMediaElement{
  captureStream(): MediaStream;
}

type KecpChat = KecpMessage & {chatID: number}

const roomID = ref('');
const username = ref('');
const conn = ref<KecpConnection>();
const userList = ref<string[]>([]);
const chatList = ref<KecpChat[]>([]);
const sig = new KecpSignal('http://127.0.0.1:8090/api/kecp/');
const connected = ref(false);
const sendValue = ref('');
const chatIDCouter = ref(0);
const chatBox = ref<HTMLElement[]>([]);
const videoInput = ref<HTMLInputElement | null>(null);
const video = ref<HTMLMediaElementWithCaputreStream | null>(null);
const fileURL = ref('');

async function connect() {
  if (roomID.value === '') {
    const room = await sig.newRoom();
    if (!room.errorText) {
      roomID.value = room.roomID;
    }
  }

  conn.value = sig.newConnection({
    roomID: roomID.value,
    name: username.value,
  });

  conn.value.on(KecpEventType.Open, () => {
    connected.value = true;
  });

  conn.value.on(KecpEventType.UserListInit, (evt: CustomEvent<string[]>) => {
    userList.value = evt.detail.slice();
  });

  conn.value.on(KecpEventType.UserJoin, (evt: CustomEvent<string>) => {
    userList.value.push(evt.detail);
  });

  conn.value.on(KecpEventType.UserLeave, (evt: CustomEvent<string>) => {
    if (userList.value.indexOf(evt.detail) !== -1) {
      userList.value.splice(userList.value.indexOf(evt.detail));
    }
  });

  conn.value.on(KecpEventType.Chat, (evt: CustomEvent<KecpMessage>) => {
    chatIDCouter.value += 1;
    chatList.value.push({
      chatID: chatIDCouter.value,
      type: evt.detail.type,
      name: evt.detail.name,
      target: evt.detail.target,
      payload: evt.detail.payload,
    });
  });

  conn.value.on(KecpEventType.VideoOffer, (evt: CustomEvent<AnswerPeer>) => {
    evt.detail.setHandleTrackEvent((event) => {
      const stream = event.streams[0];
      if (video.value) {
        if (video.value.srcObject !== stream) {
          video.value.srcObject = stream;
          video.value.autoplay = true;
        }
      }
    });
    evt.detail.answer();
  });
}

function handleSendButton() {
  if (sendValue.value === '') {
    return;
  }
  conn.value?.send(JSON.stringify({
    type: KecpMessageType.Chat,
    name: conn.value.getName(),
    payload: sendValue.value,
  }));
}

function handleKey(event: KeyboardEvent) {
  if (event.key === 'Enter') {
    handleSendButton();
  }
}

function setChatRef(el) {
  if (el) {
    chatBox.value?.push(el);
    el.scrollIntoView({ behavior: 'smooth' });
  }
}

async function dial(user: string) {
  if (user === conn.value?.getName()) {
    return;
  }

  const stream = video.value?.captureStream();
  if (stream) {
    const offer = conn.value?.newOffer(user);
    offer?.addStream(stream);
  }
}

function copyToTheClipboard() {
  navigator.clipboard.writeText(roomID.value);
}

function pasteFromTheClipboard() {
  navigator.clipboard.readText().then((value) => { roomID.value = value; });
}

function uploadFile() {
  if (videoInput.value?.files === null) {
    return;
  }
  const file = videoInput.value?.files[0];
  if (file !== undefined) {
    if (fileURL.value !== '') {
      URL.revokeObjectURL(fileURL.value);
    }
    fileURL.value = URL.createObjectURL(file);
  }
  if (video.value !== null) {
    video.value.src = fileURL.value;
  }
}

function hangUpCall() { }
</script>

<template>
  <div class="container">
    <div class="infobox">
      <p>
        This is a simple video sharing system implemented using WebSockets.
        It works by sending packets of JSON back and forth with the server.
      </p>
      <p class="disclaimer">
        This example is offered
        as-is for demonstration purposes only,
        and should not be used for any other purpose.
      </p>
      <p>
        Click a username in the user list to ask them
        to enter a one-on-one video chat with you.
      </p>
      <p>
        Enter a room id:
        <input
          id="room"
          v-model="roomID"
          type="text"
          maxlength="64"
          minlength="64"
          inputmode="text"
        >
        <input
          type="button"
          name="copy"
          value="Copy"
          @click="copyToTheClipboard"
        >
        <input
          type="button"
          name="paste"
          value="Paste"
          @click="pasteFromTheClipboard"
        >
      </p>
      <p>
        Enter a username:
        <input
          id="name"
          v-model="username"
          type="text"
          maxlength="12"
          required
          autocomplete="username"
          inputmode="text"
          placeholder="Username"
        >
        <input
          type="button"
          name="login"
          value="Log in"
          @click="connect"
        >
      </p>
    </div>
    <ul class="userlistbox">
      <li
        v-for="user in userList"
        :key="user"
        @click="dial(user)"
      >
        {{ user }}
      </li>
    </ul>
    <div class="chatbox-container">
      <div
        class="chatbox"
      >
        <span
          v-for="chat in chatList"
          :ref="setChatRef"
          :key="chat.chatID"
        >
          {{ chat.name ?? '' }}: {{ chat.payload }} <br>
        </span>
      </div>
    </div>
    <div class="videobox">
      <video
        id="video"
        ref="video"
        controls
      />

      <button
        id="hangup-button"
        role="button"
        disabled
        @click="hangUpCall"
      >
        Hang Up
      </button>
    </div>
    <div class="uploadbox">
      <input
        id="videoInput"
        ref="videoInput"
        type="file"
        :disabled="!connected"
        @change="uploadFile"
      >
    </div>
    <div class="chat-controls">
      Chat:<br>
      <input
        id="text"
        v-model="sendValue"
        type="text"
        name="text"
        size="100"
        maxlength="256"
        placeholder="Say something meaningful..."
        autocomplete="off"
        :disabled="!connected"
        @keyup="handleKey"
      >
      <input
        id="send"
        type="button"
        value="Send"
        :disabled="!connected"
        @click="handleSendButton()"
      >
    </div>
  </div>
</template>

<style scoped>

.disclaimer {
  font-size:18px;
  background-color: #ddd;
  color: black;
  margin-left: 80px;
  margin-right: 80px;
  max-width: 620px;
  padding: 12px;
  border-radius: 5px;
  border: 1px solid black;
  box-shadow: 1px 1px 2px black;
}

.container {
  display: grid;
  min-width: 1250px;
  height: 100%;
  grid-template-areas: "infobox infobox infobox"
    "userlistbox chatbox-container camerabox"
    "empty-container chat-controls chat-controls";
  grid-template-columns: 10em 1fr 500px;
  grid-template-rows: 18em 1fr 5em;
  grid-gap: 1rem;
}

.infobox {
  grid-area: infobox;
  overflow: auto;
}

.userlistbox {
  grid-area: userlistbox;
  border: 1px solid black;
  margin: 0;
  padding: 1px;
  list-style: none;
  line-height: 1.1;
  overflow-y: auto;
  overflow-x: hidden;
}

.userlistbox li {
  cursor: pointer;
  padding: 1px;
}

.chatbox-container {
  grid-area: chatbox-container;
  position: relative;
}

.chatbox {
  height: 100%;
  width: calc(100% - 18px);
  position: absolute;
  left: 0;
  top: 0;
  border: 1px solid black;
  margin: 0;
  overflow-y: scroll;
  padding: 1px;
  padding: 0.1rem 0.5rem;
}

.videobox {
  grid-area: camerabox;
  width: 500px;
  height: 375px;
  border: 1px solid black;
  display: block;
  position: relative;
  overflow: auto;
}

#received_video {
  width: 100%;
  height: 100%;
  position: absolute;
}

/* The small "preview" view of your camera */
#local_video {
  width: 120px;
  height: 90px;
  position: absolute;
  top: 1rem;
  left: 1rem;
  border: 1px solid rgba(255, 255, 255, 0.75);
  box-shadow: 0 0 4px black;
}

/* The "Hang up" button */
#hangup-button {
  display: block;
  width: 80px;
  height: 24px;
  border-radius: 8px;
  position: relative;
  margin: auto;
  top: calc(100% - 40px);
  background-color: rgba(150, 0, 0, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.7);
  box-shadow: 0px 0px 1px 2px rgba(0, 0, 0, 0.4);
  font-size: 14px;
  font-family: "Lucida Grande", "Arial", sans-serif;
  color: rgba(255, 255, 255, 1.0);
  cursor: pointer;
}

#hangup-button:hover {
  filter: brightness(150%);
  -webkit-filter: brightness(150%);
}

#hangup-button:disabled {
  filter: grayscale(50%);
  -webkit-filter: grayscale(50%);
  cursor: default;
}

.uploadbox {
  grid-area: empty-container;
}

.chat-controls {
  grid-area: chat-controls;
  width: 100%;
  height: 100%;
}
</style>
