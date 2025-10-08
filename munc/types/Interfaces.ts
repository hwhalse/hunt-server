import { UUID } from "crypto";
import { Socket } from "net";

// Munc message refers to incoming messages, which are flexible to accomodate TCP, UDP, ACK, Multi-packet, and Heartbeat
// Do not use for sending messages, use more specific types to ensure continuity
export interface MuncMessage {
    id: string;
    userId?: UUID;
    ackRequired?: boolean;
    roomName?: string;
    muncEventId: number;
    timestamp?: string;
    data: any;
    multiPacket?: {
        expectedCount: number;
        index: number;
        received: Array<number>;
    };
}

export interface MuncAckMessage {
    id: string;
    timestamp: string;
    userId: string;
    multiPacket?: {
        received: Array<number>;
    }
}

export interface HeartbeatMessage {
    id: string;
    muncEventId: number;
}

export const TCPEventIds = {
    HEARTBEAT: 512,
    INIT: 513,
    JOIN_ROOM: 514,
    ROOM_JOIN_SUCCESS: 515,
    ANNOUNCE_ROOM_JOIN: 516,
    EXIT_ROOM: 517,
    ANNOUNCE_ROOM_EXIT: 518,
    SERVER_FAILURE: 1023
}

export const PositionEventIds = {
    HEARTBEAT: 1024,
    PORT_INIT: 1025,
    CREATE_OBJECT: 1026,
    UPDATE_OBJECT: 1027,
}

export const PositionEventCallbacks = {

}

export const VoipEventIds = {
    HEARTBEAT: 2048,
    PORT_INIT: 2049,
    MESSAGE: 2050
}