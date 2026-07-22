export type PbxRegistration = {
    aor: string;
    number: string;
    contact: string;
    expires: number;
    updatedAt: string;
};

export type PbxCall = {
    id: string;
    direction: string;
    state: string;
    from: string;
    to: string;
    startedAt: string;
    answeredAt?: string;
    endedAt?: string;
    durationSec: number;
};