import { IChat } from "../organisms/chat/IChats.ts";

type EventTypes = {
  "chat_created": { roomId: string; chat: IChat };
  "chat_edited": { roomId: string; chat: IChat };
  "chat_deleted": { roomId: string; chat: IChat };
  "notification": { message: string };
};

type EventName = keyof EventTypes;

type EventCallback<K extends EventName> = (payload: EventTypes[K]) => void;

export const createEventEmitter = () => {
  const listeners: Partial<Record<EventName, Set<EventCallback<any>>>> = {};

  const on = <K extends EventName>(
    eventName: K,
    callback: EventCallback<K>,
  ) => {
    if (!listeners[eventName]) {
      listeners[eventName] = new Set();
    }
    listeners[eventName].add(callback);
  };

  const off = <K extends EventName>(
    eventName: K,
    callback: EventCallback<K>,
  ) => {
    if (listeners[eventName]) {
      listeners[eventName].delete(callback);
    }
  };

  const emit = <K extends EventName>(eventName: K, payload: EventTypes[K]) => {
    if (listeners[eventName]) {
      listeners[eventName].forEach((callback) => {
        Promise.resolve().then(() => {
          callback(payload);
        });
      });
    }
  };

  return {
    on,
    off,
    emit,
  };
};

export const EventEmitter = createEventEmitter();
