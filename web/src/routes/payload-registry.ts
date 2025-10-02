import { PlayerInfoPayload, MatchingPayload, GameStartPayload, GuessPayload, RoundStartPayload, InvalidWordPayload, GuessTimeoutPayload, FeedbackPayload, GameOverPayload, TypingPayload, PlayAgainPayload } from "$lib/types/payload";
import type { Payload } from "$lib/utils/message";

export const payloadRegistry = new Map<string, () => Payload>();
payloadRegistry.set('player_info', () => new PlayerInfoPayload());
payloadRegistry.set('matching', () => new MatchingPayload());
payloadRegistry.set('game_start', () => new GameStartPayload());
payloadRegistry.set('guess', () => new GuessPayload());
payloadRegistry.set('round_start', () => new RoundStartPayload());
payloadRegistry.set('invalid_word', () => new InvalidWordPayload());
payloadRegistry.set('guess_timeout', () => new GuessTimeoutPayload());
payloadRegistry.set('feedback', () => new FeedbackPayload());
payloadRegistry.set('game_over', () => new GameOverPayload());
payloadRegistry.set('typing', () => new TypingPayload());
payloadRegistry.set('play_again', () => new PlayAgainPayload());