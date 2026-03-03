import { useCallback, useEffect, useRef, useState } from "react";
import type {
  Board,
  Cell,
  ClientMessage,
  GameStatus,
  Player,
  ServerMessage,
  WinningCells,
} from "./types";

const EMPTY_BOARD: Board = [
  ["", "", ""],
  ["", "", ""],
  ["", "", ""],
];

type ConnectionState = "connecting" | "connected" | "disconnected";

export interface GameState {
  connectionState: ConnectionState;
  mySymbol: Player | null;
  board: Board;
  turn: Player;
  status: GameStatus;
  winner: Player | "";
  winningCells: WinningCells;
  statusMessage: string;
  errorMessage: string | null;
}

export interface GameActions {
  sendMove: (row: number, col: number) => void;
  sendRestart: () => void;
  sendQuit: () => void;
}

export function useGame(): GameState & GameActions {
  const ws = useRef<WebSocket | null>(null);

  const [connectionState, setConnectionState] =
    useState<ConnectionState>("connecting");
  const [mySymbol, setMySymbol] = useState<Player | null>(null);
  const [board, setBoard] = useState<Board>(EMPTY_BOARD);
  const [turn, setTurn] = useState<Player>("X");
  const [status, setStatus] = useState<GameStatus>("waiting");
  const [winner, setWinner] = useState<Player | "">("");
  const [winningCells, setWinningCells] = useState<WinningCells>([]);
  const [statusMessage, setStatusMessage] = useState("Connecting to server...");
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const send = useCallback((msg: ClientMessage) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(msg));
    }
  }, []);

  // Flash an error for 2.5 seconds without losing the status message.
  const flashError = useCallback((msg: string) => {
    setErrorMessage(msg);
    setTimeout(() => setErrorMessage(null), 2500);
  }, []);

  useEffect(() => {
    const socket = new WebSocket(`ws://${location.host}/ws`);
    ws.current = socket;

    socket.onopen = () => {
      setConnectionState("connected");
    };

    socket.onclose = () => {
      setConnectionState("disconnected");
      setStatusMessage("Disconnected from server.");
      setStatus("waiting");
    };

    socket.onerror = () => {
      setStatusMessage("Connection error. Is the server running?");
    };

    socket.onmessage = (event: MessageEvent<string>) => {
      const msg = JSON.parse(event.data) as ServerMessage;

      switch (msg.type) {
        case "player_assigned":
          setMySymbol(msg.player);
          setStatusMessage(msg.message);
          break;

        case "opponent_joined":
          setStatusMessage(msg.message);
          break;

        case "game_state": {
          setBoard(msg.board);
          setTurn(msg.turn);
          setStatus(msg.status);
          setWinningCells([]);
          setWinner("");
          break;
        }

        case "game_over":
          setWinner(msg.winner);
          setWinningCells(msg.winning_cells ?? []);
          setStatusMessage(msg.message);
          setStatus(msg.winner ? "won" : "draw");
          break;

        case "opponent_left":
          setStatusMessage(msg.message);
          setBoard(EMPTY_BOARD);
          setStatus("waiting");
          setWinningCells([]);
          setWinner("");
          break;

        case "error":
          flashError(msg.message);
          break;

        case "quit":
          setStatusMessage(msg.message);
          break;
      }
    };

    return () => {
      socket.close();
    };
  }, [flashError]);

  const sendMove = useCallback(
    (row: number, col: number) => send({ type: "move", row, col }),
    [send],
  );
  const sendRestart = useCallback(() => send({ type: "restart" }), [send]);
  const sendQuit = useCallback(() => send({ type: "quit" }), [send]);

  return {
    connectionState,
    mySymbol,
    board,
    turn,
    status,
    winner,
    winningCells,
    statusMessage,
    errorMessage,
    sendMove,
    sendRestart,
    sendQuit,
  };
}

export function isWinningCell(
  winningCells: WinningCells,
  row: number,
  col: number,
): boolean {
  return winningCells.some(([r, c]) => r === row && c === col);
}

export function getCellOwner(board: Board, row: number, col: number): Cell {
  return board[row][col];
}
