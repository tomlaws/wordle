import { expect, test } from 'vitest';
import Match from './Match.svelte';
import { render, waitFor } from '@testing-library/svelte';
import { createMockWebSocket } from '$lib/utils/websocket';
import type { Payload } from '$lib/utils/message';
import { InvalidWordPayload, RoundStartPayload } from '$lib/types/payload';
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
  render(Match, { props: { websocket, playerInfo }, context: new Map([[TOAST_KEY, mockToast]]) });

  const invalidWordPayload = new InvalidWordPayload();
  invalidWordPayload.player = playerInfo;
  invalidWordPayload.word = 'INVALID';
  websocket.send(invalidWordPayload);

  await waitFor(() => {
    expect(mockToast.error).toHaveBeenCalledWith(expect.stringContaining('Invalid word'));
  });
});