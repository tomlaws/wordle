import type { FeedbackPayload, GameOverPayload } from "$lib/types/payload";
import type { Payload } from "$lib/utils/message";
import type { WebSocketConnection } from "$lib/utils/websocket";

export interface GameContext {
    loading: boolean;
    websocket: WebSocketConnection<Payload>;
    playerInfo: {
        id: string; nickname: string
    };
    matchInfo: {
        player1: { id: string; nickname: string };
        player2: { id: string; nickname: string };
        guesses: Array<FeedbackPayload['feedback']>;
        currentGuess: Array<string>;
        currentRound: number;
        myTurn: boolean;
        deadline?: Date;
        gameOver?: GameOverPayload;
    } | null;
}

export const GAME_KEY = Symbol('game');