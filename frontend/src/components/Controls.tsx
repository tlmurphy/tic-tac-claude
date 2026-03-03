interface ControlsProps {
  visible: boolean;
  onRestart: () => void;
  onQuit: () => void;
}

export function Controls({ visible, onRestart, onQuit }: ControlsProps) {
  if (!visible) return null;

  return (
    <div className="flex gap-4 justify-center mt-2">
      <button
        onClick={onRestart}
        className="px-6 py-2.5 rounded-xl bg-emerald-500/20 text-emerald-300 ring-1 ring-emerald-500/40
                   hover:bg-emerald-500/30 hover:ring-emerald-400 hover:text-emerald-200
                   transition-all duration-200 font-semibold text-sm cursor-pointer"
      >
        Play Again
      </button>
      <button
        onClick={onQuit}
        className="px-6 py-2.5 rounded-xl bg-white/5 text-gray-400 ring-1 ring-white/10
                   hover:bg-rose-500/20 hover:text-rose-300 hover:ring-rose-500/40
                   transition-all duration-200 font-semibold text-sm cursor-pointer"
      >
        Quit
      </button>
    </div>
  );
}
