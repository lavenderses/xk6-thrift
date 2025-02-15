namespace java idl

service TestService {
    string simpleCall(1: string id);

    bool boolCall(1: bool tf);
}
