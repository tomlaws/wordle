import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Observable, Subscriber } from 'rxjs';

export type WebSocketConnection<T> = {
  send: (msg: T) => void;
  messages$: Observable<T>;
};

export function createWebSocket<T>(
  url: string,
  wrap: (msg: T) => any,
  unwrap: (msg: any) => T
): WebSocketConnection<T> {
  const socket$ = webSocket<any>(url);

  // Wrap sending
  const send = (msg: T) => {
    const wrapped = wrap(msg);
    socket$.next(wrapped);
  };

  // Unwrap receiving
  const messages$ = new Observable<T>((subscriber: Subscriber<T>) => {
    const sub = socket$.subscribe({
      next: raw => subscriber.next(unwrap(raw)),
      error: err => subscriber.error(err),
      complete: () => subscriber.complete()
    });

    return () => sub.unsubscribe();
  });

  return { send, messages$ };
}

export function createMockWebSocket<T>(): WebSocketConnection<T> {
  let messageSubscriber: Subscriber<T>;

  const messages$ = new Observable<T>((subscriber: Subscriber<T>) => {
    messageSubscriber = subscriber;
  });

  const send = (msg: T) => {
    if (messageSubscriber) {
      messageSubscriber.next(msg);
    }
  };

  return { send, messages$ };
}
