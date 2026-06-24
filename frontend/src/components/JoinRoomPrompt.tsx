import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Card } from "./ui/card";

interface JoinRoomPromptProps {
    onSubmit: (displayName: string) => void;
}

export function JoinRoomPrompt({
    onSubmit,
}: JoinRoomPromptProps) {
    const [displayName, setDisplayName] = useState('')

    const handleSubmit = () => {
        const trimmed = displayName.trim()

        if (!trimmed) return

        onSubmit(trimmed)
    }

    return (
        <div className="flex min-h-[calc(100vh-3.5rem)] items-center justify-center p-4">
            <Card className="w-full max-w-md space-y-4">
                <h1 className="text-xl font-semibold">
                    Join Room
                </h1>

                <p className="text-sm text-zinc-400">
                    Enter your display name to join this room.
                </p>

                <Input
                    id="display-name"
                    label="Display Name"
                    placeholder="Your name"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            handleSubmit()
                        }
                    }}
                />

                <Button
                    className="w-full"
                    onClick={handleSubmit}
                >
                    Join Room
                </Button>
            </Card>
        </div>
    )
}