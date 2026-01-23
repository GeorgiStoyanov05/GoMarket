"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Star, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";

export default function WatchlistButton({
  symbol,
  company,
  isInWatchlist: initialIsInWatchlist,
  showTrashIcon = false,
  type = "button",
  onWatchlistChange,
}: WatchlistButtonProps) {
  const [isInWatchlist, setIsInWatchlist] = useState(initialIsInWatchlist);
  const [isLoading, setIsLoading] = useState(false);

  const handleToggle = async () => {
    setIsLoading(true);
    try {
      // TODO: Implement API call to add/remove from watchlist
      // const response = await fetch(`/api/watchlist/${symbol}`, {
      //   method: isInWatchlist ? 'DELETE' : 'POST',
      //   body: JSON.stringify({ symbol, company }),
      // });
      
      // For now, just toggle the local state
      const newState = !isInWatchlist;
      setIsInWatchlist(newState);
      
      // Call the callback if provided
      if (onWatchlistChange) {
        onWatchlistChange(symbol, newState);
      }
    } catch (error) {
      console.error("Error toggling watchlist:", error);
      // Revert on error
      setIsInWatchlist(isInWatchlist);
    } finally {
      setIsLoading(false);
    }
  };

  if (type === "icon") {
    return (
      <button
        onClick={handleToggle}
        disabled={isLoading}
        className={cn(
          "watchlist-icon-btn",
          isInWatchlist && "watchlist-icon-added"
        )}
        aria-label={isInWatchlist ? "Remove from watchlist" : "Add to watchlist"}
      >
        <div className="watchlist-icon">
          {showTrashIcon && isInWatchlist ? (
            <Trash2 className="trash-icon" />
          ) : (
            <Star className={cn("star-icon", isInWatchlist && "fill-current")} />
          )}
        </div>
      </button>
    );
  }

  return (
    <Button
      onClick={handleToggle}
      disabled={isLoading}
      className={cn(
        "watchlist-btn",
        isInWatchlist && "watchlist-remove"
      )}
    >
      {isInWatchlist ? "Remove from Watchlist" : "Add to Watchlist"}
    </Button>
  );
}
