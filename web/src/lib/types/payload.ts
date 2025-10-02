export class PlayerInfoPayload {
    id!: string;
    nickname!: string;
    
    MessageType(): string {
        return 'player_info';
    }
}

export class MatchingPayload{
    MessageType(): string {
        return 'matching';
    }
}

export class GameStartPayload {
    maxGuesses!: number;
    player1!: { id: string; nickname: string; };
    player2!: { id: string; nickname: string; };

    MessageType(): string {
        return 'game_start';
    }
}

export class GuessPayload {
    word!: string;

    MessageType(): string {
        return 'guess';
    }
}

export class RoundStartPayload {
    player!: { id: string; nickname: string; };
    round!: number;
    deadline!: string;

    getDeadline(): Date {
        return new Date(this.deadline);
    }

    MessageType(): string {
        return 'round_start';
    }
}

export class InvalidWordPayload {
    player!: { id: string; nickname: string; };
    round!: number;
    word!: string;

    MessageType(): string {
        return 'invalid_word';
    }
}

export class GuessTimeoutPayload {
    player!: { id: string; nickname: string; };
    round!: number;

    MessageType(): string {
        return 'guess_timeout';
    }
}

export class FeedbackPayload {
    player!: { id: string; nickname: string; };
    round!: number;
    feedback!: Array<{
        letter: number;
        position: number;
        matchType: number;
    }>;

    MessageType(): string {
        return 'feedback';
    }
}

export class GameOverPayload {
    winner!: { id: string; nickname: string; } | null;
    answer!: string;

    MessageType(): string {
        return 'game_over';
    }
}

export class TypingPayload {
    player?: { id: string; nickname: string; };
    word!: string;

    MessageType(): string {
        return 'typing';
    }
}

export class PlayAgainPayload {
    confirm!: boolean;

    MessageType(): string {
        return 'play_again';
    }
}