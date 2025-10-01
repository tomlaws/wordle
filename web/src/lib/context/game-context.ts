import type { Payload } from "$lib/utils/message";
import type { WebSocketConnection } from "$lib/utils/websocket";

export interface GameContext {
    websocket: WebSocketConnection<Payload>;
    playerInfo: { id: string; nickname: string };
    matchInfo: { player1: { id: string; nickname: string }; player2: { id: string; nickname: string } } | null;
}

export const GAME_KEY = Symbol('game');