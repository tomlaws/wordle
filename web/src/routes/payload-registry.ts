import { FeedbackPayload, GuessPayload, InvalidWordPayload, PlayerInfoPayload, RoundStartPayload } from "$lib/types/payload";
import type { Payload } from "$lib/utils/message";

export const payloadRegistry = new Map<string, () => Payload>();
payloadRegistry.set('player_info', () => new PlayerInfoPayload());
payloadRegistry.set('guess', () => new GuessPayload());
payloadRegistry.set('round_start', () => new RoundStartPayload());
payloadRegistry.set('invalid_word', () => new InvalidWordPayload());
payloadRegistry.set('feedback', () => new FeedbackPayload());