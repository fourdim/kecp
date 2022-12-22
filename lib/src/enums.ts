/* eslint-disable no-unused-vars */
export enum KecpMessageType {
    VideoOffer = 'video-offer',
    VideoAnswer = 'video-answer',
    DataOffer = 'data-offer',
    DataAnswer = 'data-answer',
    NewIceCandidate = 'new-ice-candidate',
    Chat = 'chat',
    List = 'list',
    Join = 'join',
    Leave = 'leave',
    Error = 'error',
}

export enum KecpEventType {
    Error = 'error',
    VideoOffer = 'video-offer',
    VideoAnswer = 'video-answer',
    DataOffer = 'data-offer',
    DataAnswer = 'data-answer',
    NewIceCandidate = 'new-ice-candidate',
    Chat = 'chat',
    Open = 'open',
    UserListInit = 'userlist-init',
    UserJoin = 'userlist-join',
    UserLeave = 'userlist-leave',
}
