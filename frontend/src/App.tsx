import { useGame } from "./useGame";
import { Board } from "./components/Board";
import { StatusBar } from "./components/StatusBar";
import { Controls } from "./components/Controls";

export default function App() {
  const {
    connectionState,
    mySymbol,
    board,
    turn,
    status,
    winner: _winner,
    winningCells,
    statusMessage,
    errorMessage,
    sendMove,
    sendRestart,
    sendQuit,
  } = useGame();

  const canPlay = status === "playing";
  const gameOver = status === "won" || status === "draw";

  return (
    <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center p-4">
      <div className="flex flex-col items-center gap-8 w-full max-w-sm">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-3xl font-bold tracking-tight text-white">
            Tic-Tac-Toe
          </h1>
          <p className="text-gray-500 text-xs mt-1 uppercase tracking-widest">
            Local Multiplayer
          </p>
        </div>

        {/* Status */}
        <StatusBar
          connectionState={connectionState}
          mySymbol={mySymbol}
          turn={turn}
          status={status}
          statusMessage={statusMessage}
          errorMessage={errorMessage}
        />

        {/* Board */}
        <div
          className={`transition-opacity duration-300 ${
            status === "waiting" || connectionState !== "connected"
              ? "opacity-40 pointer-events-none"
              : "opacity-100"
          }`}
        >
          <Board
            board={board}
            mySymbol={mySymbol}
            turn={turn}
            canPlay={canPlay}
            winningCells={winningCells}
            onCellClick={sendMove}
          />
        </div>

        {/* Post-game controls */}
        <Controls
          visible={gameOver && connectionState === "connected"}
          onRestart={sendRestart}
          onQuit={sendQuit}
        />

        {/* Waiting indicator */}
        {status === "waiting" && connectionState === "connected" && (
          <div className="flex items-center gap-2 text-gray-500 text-sm">
            <span className="inline-block w-2 h-2 rounded-full bg-amber-400 animate-pulse" />
            Waiting for opponent...
          </div>
        )}
      </div>
    </div>
  );
}
