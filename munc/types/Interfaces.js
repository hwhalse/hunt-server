"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.VoipEventIds = exports.PositionEventCallbacks = exports.PositionEventIds = exports.TCPEventIds = void 0;
exports.TCPEventIds = {
    HEARTBEAT: 512,
    INIT: 513,
    JOIN_ROOM: 514,
    ROOM_JOIN_SUCCESS: 515,
    ANNOUNCE_ROOM_JOIN: 516,
    EXIT_ROOM: 517,
    ANNOUNCE_ROOM_EXIT: 518,
    SERVER_FAILURE: 1023
};
exports.PositionEventIds = {
    HEARTBEAT: 1024,
    PORT_INIT: 1025,
    CREATE_OBJECT: 1026,
    UPDATE_OBJECT: 1027,
};
exports.PositionEventCallbacks = {};
exports.VoipEventIds = {
    HEARTBEAT: 2048,
    PORT_INIT: 2049,
    MESSAGE: 2050
};
