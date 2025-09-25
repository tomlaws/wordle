import type { Payload } from "$lib/utils/message";
import type { WebSocketConnection } from "$lib/utils/websocket";

export interface GameContext {
    websocket: WebSocketConnection<Payload>;
    playerInfo: { id: string; nickname: string };
}

export const GAME_KEY = Symbol('game');