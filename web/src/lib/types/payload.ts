export class PlayerInfoPayload {
    id!: string;
    nickname!: string;
    MessageType(): string {
        return 'player_info';
    }
}

export class GuessPayload {
    guess!: string;

    constructor();
    constructor(init: Partial<GuessPayload>);
    constructor(init?: Partial<GuessPayload>) {
        Object.assign(this, init);
    }

    MessageType(): string {
        return 'guess';
    }
}

export class RoundStartPayload {
    player!: { id: string; nickname: string; };
    round!: number;
    timeout!: number;
    MessageType(): string {
        return 'round_start';
    }
}

export class InvalidWordPayload {
    player!: { id: string; nickname: string; };
    word!: string;

    constructor();
    constructor(init: Partial<InvalidWordPayload>);
    constructor(init?: Partial<InvalidWordPayload>) {
        Object.assign(this, init);
    }

    MessageType(): string {
        return 'invalid_word';
    }
}

export class FeedbackPayload {
    player!: { id: string; nickname: string; };
    round!: number;
    feedback!: Array<{
        letter: string;
        position: number;
        matchType: number;
    }>;

    constructor();
    constructor(init: Partial<FeedbackPayload>);
    constructor(init?: Partial<FeedbackPayload>) {
        Object.assign(this, init);
    }

    MessageType(): string {
        return 'feedback';
    }
}