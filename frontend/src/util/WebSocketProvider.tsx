import { createContext, useContext, useEffect, useRef, useState } from "react";

export const WebSocketContext = createContext<WebSocket | null>(null);

const normalizeWsUrl = (value: string | undefined) => {
  if (!value || value.trim() === "") return "/ws";
  return value.trim();
};

const wsBase = normalizeWsUrl(import.meta.env.VITE_REACT_APP_WS_BASE_URL);

const resolveWsUrl = () => {
  if (wsBase.startsWith("ws://") || wsBase.startsWith("wss://")) {
    return wsBase;
  }
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const path = wsBase.startsWith("/") ? wsBase : `/${wsBase}`;
  return `${protocol}//${window.location.host}${path}`;
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("WebSocketContext is not provided");
  }
  return context;
};

export const WebSocketProvider = (
  { children }: { children: React.ReactNode },
) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const connectionAttemptedRef = useRef(false);

  useEffect(() => {
    if (connectionAttemptedRef.current) return;
    connectionAttemptedRef.current = true;

    if (
      socketRef.current && socketRef.current.readyState !== WebSocket.CLOSED
    ) {
      console.log("Closing existing WebSocket connection");
      socketRef.current.close();
      socketRef.current = null;
      setSocket(null);
    }

    const ws = new WebSocket(resolveWsUrl());
    ws.addEventListener("open", () => {
      console.log("WebSocket connection established");
      socketRef.current = ws;
      setSocket(ws);
    });

    ws.addEventListener("error", (event) => {
      console.error("WebSocket error observed:", event);
      connectionAttemptedRef.current = false;
    });

    ws.addEventListener("close", () => {
      console.log("closed WebSocket");
      socketRef.current = null;
      setSocket(null);
      connectionAttemptedRef.current = false;
    });

    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, []);

  if (!socket) {
    return <></>;
  }

  return (
    <WebSocketContext.Provider value={socket}>
      {children}
    </WebSocketContext.Provider>
  );
};
