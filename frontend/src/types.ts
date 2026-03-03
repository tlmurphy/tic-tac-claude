export type Player = "X" | "O";
export type Cell = "X" | "O" | "";
export type Board = [
  [Cell, Cell, Cell],
  [Cell, Cell, Cell],
  [Cell, Cell, Cell],
];
export type GameStatus = "waiting" | "playing" | "won" | "draw";
export type WinningCells = [number, number][];

// --- Server → Client messages ---

export interface MsgPlayerAssigned {
  type: "player_assigned";
  player: Player;
  message: string;
}

export interface MsgOpponentJoined {
  type: "opponent_joined";
  message: string;
}

export interface MsgGameState {
  type: "game_state";
  board: Board;
  turn: Player;
  status: GameStatus;
}

export interface MsgGameOver {
  type: "game_over";
  winner: Player | "";
  winning_cells: WinningCells;
  message: string;
}

export interface MsgOpponentLeft {
  type: "opponent_left";
  message: string;
}

export interface MsgError {
  type: "error";
  message: string;
}

export interface MsgQuit {
  type: "quit";
  message: string;
}

export type ServerMessage =
  | MsgPlayerAssigned
  | MsgOpponentJoined
  | MsgGameState
  | MsgGameOver
  | MsgOpponentLeft
  | MsgError
  | MsgQuit;

// --- Client → Server messages ---

export interface OutMove {
  type: "move";
  row: number;
  col: number;
}

export interface OutRestart {
  type: "restart";
}

export interface OutQuit {
  type: "quit";
}

export type ClientMessage = OutMove | OutRestart | OutQuit;
