import type { GameStatus, Player } from "../types";

interface StatusBarProps {
  connectionState: "connecting" | "connected" | "disconnected";
  mySymbol: Player | null;
  turn: Player;
  status: GameStatus;
  statusMessage: string;
  errorMessage: string | null;
}

export function StatusBar({
  connectionState,
  mySymbol,
  turn,
  status,
  statusMessage,
  errorMessage,
}: StatusBarProps) {
  if (connectionState !== "connected") {
    return (
      <div className="text-center">
        <p className="text-gray-400 text-sm">{statusMessage}</p>
      </div>
    );
  }

  const isMyTurn = status === "playing" && turn === mySymbol;

  return (
    <div className="text-center space-y-2">
      {/* Player badge */}
      {mySymbol && (
        <div className="flex items-center justify-center gap-2">
          <span className="text-gray-400 text-sm">You are</span>
          <span
            className={`text-lg font-bold px-3 py-0.5 rounded-full ${
              mySymbol === "X"
                ? "text-blue-400 bg-blue-500/10 ring-1 ring-blue-500/30"
                : "text-rose-400 bg-rose-500/10 ring-1 ring-rose-500/30"
            }`}
          >
            {mySymbol}
          </span>
        </div>
      )}

      {/* Turn / status message */}
      {errorMessage ? (
        <p className="text-rose-400 text-sm font-medium">{errorMessage}</p>
      ) : (
        <p
          className={`text-sm font-medium ${
            isMyTurn
              ? "text-emerald-400"
              : status === "won" || status === "draw"
                ? "text-yellow-300"
                : "text-gray-300"
          }`}
        >
          {isMyTurn
            ? "Your turn!"
            : status === "playing"
              ? `Waiting for ${turn} to move...`
              : statusMessage}
        </p>
      )}

      {/* Turn indicator dots */}
      {status === "playing" && (
        <div className="flex items-center justify-center gap-2 mt-1">
          <span
            className={`text-xs px-2 py-0.5 rounded-full ${
              turn === "X"
                ? "bg-blue-500/20 text-blue-300 ring-1 ring-blue-400"
                : "bg-white/5 text-gray-500"
            }`}
          >
            X
          </span>
          <span className="text-gray-600">·</span>
          <span
            className={`text-xs px-2 py-0.5 rounded-full ${
              turn === "O"
                ? "bg-rose-500/20 text-rose-300 ring-1 ring-rose-400"
                : "bg-white/5 text-gray-500"
            }`}
          >
            O
          </span>
        </div>
      )}
    </div>
  );
}
