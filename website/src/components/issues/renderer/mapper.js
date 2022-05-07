export function mapPickStatusToBackend(pick) {
    return {
        unknown: "UnKnown",
        approved: "Accept",
        later: "Later",
        "won't fix": "Won't Fix",
        "approve(frozen)": "Accept(Frozen)"
    }[pick]
}

export function mapPickStatusToFrontend(pick) {
    pick = pick.toLocaleLowerCase();
    pick = {
        accept: "approved",
        unknown: "unknown",
        later: "later",
        "won't fix": "won't fix",
        "accept(frozen)": "approved(frozen)"
    }[pick]
    return pick
}
