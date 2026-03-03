import type { Board as BoardType, Player, WinningCells } from "../types";
import { getCellOwner, isWinningCell } from "../useGame";

interface BoardProps {
  board: BoardType;
  mySymbol: Player | null;
  turn: Player;
  canPlay: boolean;
  winningCells: WinningCells;
  onCellClick: (row: number, col: number) => void;
}

export function Board({
  board,
  mySymbol,
  turn,
  canPlay,
  winningCells,
  onCellClick,
}: BoardProps) {
  return (
    <div className="grid grid-cols-3 gap-3 p-2">
      {[0, 1, 2].map((row) =>
        [0, 1, 2].map((col) => {
          const value = getCellOwner(board, row, col);
          const winning = isWinningCell(winningCells, row, col);
          const isEmpty = value === "";
          const isMyTurn = canPlay && isEmpty && turn === mySymbol;

          return (
            <Cell
              key={`${row}-${col}`}
              value={value}
              winning={winning}
              clickable={isMyTurn}
              onClick={() => isEmpty && isMyTurn && onCellClick(row, col)}
            />
          );
        }),
      )}
    </div>
  );
}

interface CellProps {
  value: string;
  winning: boolean;
  clickable: boolean;
  onClick: () => void;
}

function Cell({ value, winning, clickable, onClick }: CellProps) {
  const baseClasses =
    "flex items-center justify-center w-28 h-28 rounded-2xl text-5xl font-bold transition-all duration-200 select-none";

  const stateClasses = winning
    ? "bg-emerald-500/30 ring-2 ring-emerald-400 scale-105"
    : value === "X"
      ? "bg-blue-500/10 ring-1 ring-blue-500/30"
      : value === "O"
        ? "bg-rose-500/10 ring-1 ring-rose-500/30"
        : clickable
          ? "bg-white/5 ring-1 ring-white/10 hover:bg-white/10 hover:ring-white/20 cursor-pointer hover:scale-105"
          : "bg-white/5 ring-1 ring-white/10";

  const textColor =
    value === "X" ? "text-blue-400" : value === "O" ? "text-rose-400" : "";

  return (
    <div className={`${baseClasses} ${stateClasses}`} onClick={onClick}>
      <span className={textColor}>{value}</span>
    </div>
  );
}
