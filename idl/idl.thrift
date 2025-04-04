namespace java idl

struct Nested {
    1: string inner,
}

struct Message {
    1: string content,
    2: map<string, bool> tags,
    3: Nested nested,
}

enum Feature {
    ONE = 1;
    TWO = 2;
    THREE = 3;
}

service TestService {
    string simpleCall(1: string id);

    bool boolCall(1: bool tf);

    Message messageCall(1: Message message);

    map<string, bool> mapCall(1: map<string, bool> maps);

    list<string> stringCall(1: list<string> strs);

    list<Message> stringsCall(1: list<Message> strs);

    list<Feature> enumCall(1: Feature feature);
}
