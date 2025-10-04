import { expect, test } from 'vitest';
import Match from './Match.svelte';
import { render, screen, waitFor } from '@testing-library/svelte';
import { createMockWebSocket } from '$lib/utils/websocket';
import type { Payload } from '$lib/utils/message';
import { FeedbackPayload, GameOverPayload, GuessTimeoutPayload, InvalidWordPayload, RoundStartPayload } from '$lib/types/payload';
import { flushSync, tick } from 'svelte';
import { TOAST_KEY } from '$lib/context/toast-context';
import { mockToast } from '$lib/mocks/toast-mock';
import { GAME_KEY, type GameContext } from '$lib/context/game-context';

test('should notify when it is player\'s turn', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const gameContext = $state<GameContext>({
    websocket: websocket,
    playerInfo: playerInfo,
    matchInfo: {
      loading: true,
      player1: playerInfo,
      player2: { id: '2', nickname: 'Player2' },
      guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
      currentRound: -1,
      currentGuess: Array(5).fill(''),
      myTurn: false
    }
  });
  render(Match, {
    context: new Map<any, any>([
      [TOAST_KEY, mockToast],
      [GAME_KEY, gameContext]
    ])
  });
  const roundStartPayload = new RoundStartPayload();
  roundStartPayload.player = playerInfo;
  roundStartPayload.round = 1;
  roundStartPayload.deadline = (Date.now() + 60000).toString(); // 1 minute from now
  websocket.send(roundStartPayload);
  flushSync();
  expect(screen.getByText('Round 1')).toBeInTheDocument();
});

test('should toast invalid word message', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const gameContext = $state<GameContext>({
    websocket: websocket,
    playerInfo: playerInfo,
    matchInfo: {
      loading: true,
      player1: playerInfo,
      player2: { id: '2', nickname: 'Player2' },
      guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
      currentRound: 1,
      currentGuess: Array(5).fill(''),
      myTurn: false
    }
  });
  render(Match, {
    context: new Map<Symbol, any>([
      [TOAST_KEY, mockToast],
      [GAME_KEY, gameContext]
    ])
  });
  const invalidWordPayload = new InvalidWordPayload();
  invalidWordPayload.player = playerInfo;
  invalidWordPayload.round = 1;
  invalidWordPayload.word = 'WRONG';
  websocket.send(invalidWordPayload);
  flushSync();
  expect(mockToast.error).toHaveBeenCalledWith(expect.stringContaining('WRONG is not a valid word'));
});

test('should render correct feedback classes for each letter after guess', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const gameContext = $state<GameContext>({
    websocket: websocket,
    playerInfo: playerInfo,
    matchInfo: {
      loading: true,
      player1: playerInfo,
      player2: { id: '2', nickname: 'Player2' },
      guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
      currentRound: 1,
      currentGuess: Array(5).fill(''),
      myTurn: false
    }
  });
  const { container } = render(Match, {
    context: new Map<Symbol, any>([
      [TOAST_KEY, mockToast],
      [GAME_KEY, gameContext]
    ])
  });


  // feedback payload with all letters as 'miss'
  const feedbackPayload = new FeedbackPayload();
  feedbackPayload.player = playerInfo;
  feedbackPayload.round = 1;
  feedbackPayload.feedback = [
    { letter: 'A'.charCodeAt(0), position: 0, matchType: 0 },
    { letter: 'B'.charCodeAt(0), position: 1, matchType: 1 },
    { letter: 'C'.charCodeAt(0), position: 2, matchType: 2 },
    { letter: 'D'.charCodeAt(0), position: 3, matchType: 0 },
    { letter: 'E'.charCodeAt(0), position: 4, matchType: 1 }
  ];
  websocket.send(feedbackPayload);
  flushSync();
  const boxes = container.querySelectorAll('.row:nth-child(1) .box');
  expect(boxes[0]).toHaveClass('miss');
  expect(boxes[1]).toHaveClass('present');
  expect(boxes[2]).toHaveClass('hit');
  expect(boxes[3]).toHaveClass('miss');
  expect(boxes[4]).toHaveClass('present');
});

test('should notify when player guess times out', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const gameContext = $state<GameContext>({
    websocket: websocket,
    playerInfo: playerInfo,
    matchInfo: {
      loading: true,
      player1: playerInfo,
      player2: { id: '2', nickname: 'Player2' },
      guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
      currentRound: 1,
      currentGuess: Array(5).fill(''),
      myTurn: false
    }
  });
  const { container } = render(Match, {
    context: new Map<Symbol, any>([
      [TOAST_KEY, mockToast],
      [GAME_KEY, gameContext]
    ])
  });

  const guessTimeoutPayload = new GuessTimeoutPayload();
  guessTimeoutPayload.player = playerInfo;
  guessTimeoutPayload.round = 1;
  websocket.send(guessTimeoutPayload);
  flushSync();
  expect(mockToast.error).toHaveBeenCalledWith(expect.stringContaining('You ran out of time!'));
  const boxes = container.querySelectorAll('.row:nth-child(1) .box');
  boxes.forEach(box => {
    expect(box).toHaveClass('miss');
  });
});

test('should notify when game is over - You won', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const gameContext = $state<GameContext>({
    websocket: websocket,
    playerInfo: playerInfo,
    matchInfo: {
      loading: true,
      player1: playerInfo,
      player2: { id: '2', nickname: 'Player2' },
      guesses: Array.from({ length: 12 }, () => Array(5).fill(null)),
      currentRound: 1,
      currentGuess: Array(5).fill(''),
      myTurn: false
    }
  });
  const { container } = render(Match, {
    context: new Map<Symbol, any>([
      [TOAST_KEY, mockToast],
      [GAME_KEY, gameContext]
    ])
  });

  const gameOverPayload = new GameOverPayload();
  gameOverPayload.winner = playerInfo;
  gameOverPayload.answer = 'APPLE';
  websocket.send(gameOverPayload);
  flushSync();
  expect(container).toHaveTextContent('Congratulations');
  expect(container).toHaveTextContent('The word was APPLE');
});