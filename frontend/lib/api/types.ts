export interface IBaseContext {
    createdAt: string;
    updatedAt: string;
    deletedAt?: string;
}

export interface ITile extends IBaseContext {
    id: string;
    text: string;
    category: string;
    weight: number;
    lastDrawnShow: string;
    isActive: string;
}

export interface IShow extends IBaseContext {
    id: string;
    ytId: string;
    scheduledStart: string;
    actualStart: string;
    title: string;
    isLive: boolean;
    metadata: {
        [key: string]: unknown;
    };
}

export interface IPlayer extends IBaseContext {
    id: string;
    did: string;
    displayName: string;
    permissions: number;
    settings: {
        [key: string]: unknown;
    };
}


export interface ILeaderboardEntry {
    playerId: string;
    score: number;
}

export interface ITileConfirmation {
    id: string;
    showId: string;
    tileId: string;
    context: string;
    confirmedBy: string;
    confirmedAt: string;
}

export interface IHostLock extends IBaseContext {
    tileId: string;
    showId: string;
    lockedBy: string;
    expiresAt: string;
}

export enum IChatMessageType {
    USER = 'user',
    SYSTEM = 'system',
}

export interface IChatMessage extends IBaseContext {
    id: string;
    showId: string;
    type: IChatMessageType;
    playerId?: string;
    content: string;
}