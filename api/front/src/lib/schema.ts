export type ResponseSchema = {
    data: any,
    errors: ErrorSchema[]
}

export type ErrorSchema = {
    code: number,
    desc: string,
}

type Pair<T> = {
    first: T,
    second: T,
}

export type FriendsResponse = {
    first_id: number,
    second_id: number,
    games_in_same_lobby: number,
    avg_places: Pair<number>,
    same_lobby_avg_places: Pair<number>,
    pts_gained: Pair<number>,
}