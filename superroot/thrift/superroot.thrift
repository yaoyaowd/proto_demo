struct SearchRequest {
    1: optional string query,
    2: optional i32 offset=0,
    3: optional i32 limit=25,
}

struct SearchDoc {
    1: optional string id,
    2: optional double score,
}

struct SearchResponse {
    1: optional list<SearchDoc> docs
}

service SuperRoot {
    SearchResponse search(1:SearchRequest request),
}
