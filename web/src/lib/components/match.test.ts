import { expect, test } from 'vitest';
import Match from './Match.svelte';
import { render, waitFor } from '@testing-library/svelte';
import { createMockWebSocket } from '$lib/utils/websocket';
import type { Payload } from '$lib/utils/message';
import { FeedbackPayload, GameOverPayload, GuessTimeoutPayload, InvalidWordPayload, RoundStartPayload } from '$lib/types/payload';
import { tick } from 'svelte';
import { TOAST_KEY } from '$lib/context/toast-context';
import { mockToast } from '$lib/mocks/toast-mock';

test('should notify when it is player\'s turn', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };

  const { getByText, container } = render(Match, { props: { websocket, playerInfo } });
  await expect(getByText('Loading')).toBeInTheDocument()

  const roundStartPayload = new RoundStartPayload();
  roundStartPayload.player = playerInfo;
  roundStartPayload.round = 1;
  roundStartPayload.timeout = 60;
  websocket.send(roundStartPayload);
  await tick();
  await expect(getByText('Your Turn!')).toBeInTheDocument();
});

test('should toast invalid word message', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const { container } = render(Match, { props: { websocket, playerInfo }, context: new Map([[TOAST_KEY, mockToast]]) });

  const invalidWordPayload = new InvalidWordPayload();
  invalidWordPayload.player = playerInfo;
  invalidWordPayload.word = 'INVALID';
  websocket.send(invalidWordPayload);

  await waitFor(() => {
    expect(mockToast.error).toHaveBeenCalledWith(expect.stringContaining('Invalid word'));
  });
});

test('should render correct feedback classes for each letter after guess', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const { container } = render(Match, { props: { websocket, playerInfo }, context: new Map([[TOAST_KEY, mockToast]]) });

  // feedback payload with all letters as 'miss'
  const feedbackPayload = new FeedbackPayload();
  feedbackPayload.player = playerInfo;
  feedbackPayload.round = 1;
  feedbackPayload.feedback = [
    { letter: 'A', position: 0, matchType: 0 },
    { letter: 'B', position: 1, matchType: 1 },
    { letter: 'C', position: 2, matchType: 2 },
    { letter: 'D', position: 3, matchType: 0 },
    { letter: 'E', position: 4, matchType: 1 }
  ];
  websocket.send(feedbackPayload);

  // expect nth row will have class 'miss'
  await waitFor(() => {
    const boxes = container.querySelectorAll('.row:nth-child(1) .box');
    expect(boxes[0]).toHaveClass('miss');
    expect(boxes[1]).toHaveClass('present');
    expect(boxes[2]).toHaveClass('hit');
    expect(boxes[3]).toHaveClass('miss');
    expect(boxes[4]).toHaveClass('present');
  });
});

test('should notify when player guess times out', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const { container } = render(Match, { props: { websocket, playerInfo }, context: new Map([[TOAST_KEY, mockToast]]) });

  const guessTimeoutPayload = new GuessTimeoutPayload();
  guessTimeoutPayload.player = playerInfo;
  guessTimeoutPayload.round = 1;
  websocket.send(guessTimeoutPayload);

  await waitFor(() => {
    expect(mockToast.error).toHaveBeenCalledWith(expect.stringContaining('You ran out of time!'));
     const boxes = container.querySelectorAll('.row:nth-child(1) .box');
     boxes.forEach(box => {
       expect(box).toHaveClass('miss');
     });
  });
});

test('should notify when game is over', async () => {
  const websocket = createMockWebSocket<Payload>();
  const playerInfo = {
    id: '1',
    nickname: 'Player1'
  };
  const { container } = render(Match, { props: { websocket, playerInfo }, context: new Map([[TOAST_KEY, mockToast]]) });

  const gameOverPayload = new GameOverPayload();
  gameOverPayload.winner = playerInfo;
  gameOverPayload.answer = 'APPLE';
  websocket.send(gameOverPayload);

  await waitFor(() => {
    expect(mockToast.success).toHaveBeenCalledWith(expect.stringContaining('You won!'));
  });

  // Draw
  const gameOverDrawPayload = new GameOverPayload();
  gameOverDrawPayload.winner = null;
  gameOverDrawPayload.answer = 'APPLE';
  websocket.send(gameOverDrawPayload);

  await waitFor(() => {
    expect(mockToast.info).toHaveBeenCalledWith(expect.stringContaining('The game ended in a draw'));
  });
});