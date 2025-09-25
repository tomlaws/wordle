abstract class BasePayload<T> {
    constructor();
    constructor(init: Partial<T>);
    constructor(init?: Partial<T>) {
        Object.assign(this, init);
    }
    abstract MessageType(): string;
}

export class PlayerInfoPayload extends BasePayload<PlayerInfoPayload> {
    id!: string;
    nickname!: string;
    MessageType(): string {
        return 'player_info';
    }
}

export class GuessPayload extends BasePayload<GuessPayload> {
    guess!: string;

    MessageType(): string {
        return 'guess';
    }
}

export class RoundStartPayload extends BasePayload<RoundStartPayload> {
    player!: { id: string; nickname: string; };
    round!: number;
    timeout!: number;

    MessageType(): string {
        return 'round_start';
    }
}

export class InvalidWordPayload extends BasePayload<InvalidWordPayload> {
    player!: { id: string; nickname: string; };
    word!: string;

    MessageType(): string {
        return 'invalid_word';
    }
}

export class FeedbackPayload extends BasePayload<FeedbackPayload> {
    player!: { id: string; nickname: string; };
    round!: number;
    feedback!: Array<{
        letter: string;
        position: number;
        matchType: number;
    }>;

    MessageType(): string {
        return 'feedback';
    }
}

export class GuessTimeoutPayload extends BasePayload<GuessTimeoutPayload> {
    player!: { id: string; nickname: string; };
    round!: number;

    MessageType(): string {
        return 'guess_timeout';
    }
}