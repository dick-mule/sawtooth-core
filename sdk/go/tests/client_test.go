package tests

import (
    "testing"
    . "sawtooth_sdk/client"
)

var data = []byte{0x01, 0x02, 0x03}
var privstr = "ad8523ac9f1e7a9fdaa42c25ca766b7b099c871e7c0705ae191e0bef22b5d8cb"

func TestSigning(t *testing.T) {
    priv := GenPrivKey()
    pub := GenPubKey(priv)
    sig := Sign(data, priv)
    if !Verify(data, sig, pub) {
        t.Error(
            "Couldn't verify generated signature",
            priv, pub, sig,
        )
    }
}

func TestEncoding(t *testing.T) {
    priv := Decode(privstr, "hex")
    if Encode(priv, "hex") != privstr {
        t.Error("Private key is different after encoding/decoding")
    }

    pubstr := Encode(GenPubKey(priv), "hex")
    if pubstr != "03d6d8ab906a0ad263628e9c81f01a73dc0361c51c90a1b583aece8103126bf40c" {
        t.Error("Public key doesn't match expected. Got", pubstr)
    }
    pub := Decode(pubstr, "hex")
    if len(pub) != 33 {
        t.Error("Encoded pubkey wrong length. Should be 33, but is", len(pub))
    }

    sigstr := Encode(Sign(data, priv), "hex")
    if sigstr != "ccecded22fb1153d2f45aaf6df8280d3296dd4677b1f711d6a89fc1e81393dda5287d81569a6c91b90389c562e60169d036f1f6a66241156c69e364434bcd654" {
        t.Error("Signature doesn't match expected. Got", sigstr)
    }
}

func TestEncoder(t *testing.T) {
    priv := Decode(privstr, "hex")

    encoder := NewEncoder(priv, TransactionParams{
        FamilyName: "abc",
        FamilyVersion: "123",
        PayloadEncoding: "myencoding",
        Inputs: []string{"def"},
    })

    txn1 := encoder.NewTransaction(data, TransactionParams{
        Nonce: "123",
        Outputs: []string{"def"},
    })

    pubstr := Encode(GenPubKey(priv), "hex")
    txn2 := encoder.NewTransaction(data, TransactionParams{
        Nonce: "456",
        Outputs: []string{"ghi"},
        BatcherPubkey: pubstr,
    })

    // Test serialization
    txns, err := ParseTransactions(SerializeTransactions([]*Transaction{txn1, txn2}))
    if err != nil {
        t.Error(err)
    }

    batch := encoder.NewBatch(txns)

    // Test serialization
    batches, err := ParseBatches(SerializeBatches([]*Batch{batch}))
    if err != nil {
        t.Error(err)
    }
    data := SerializeBatches(batches)
    datastr := Encode(data, "hex")

    expected := "0acc0a0aca020a423033643664386162393036613061643236333632386539633831663031613733646330333631633531633930613162353833616563653831303331323662663430631280016231363530333032636533356330653862633432373537663134653730383633646263646539386331363734616361346637623130383531643538353863633136303965303162376130316639323536363736613939653064623035616231393435313534613232656465636466326236353038663032666564346233646362128001626436303930326430313632393863346532633238306665383061373663623966346162613936666432353933633166326335663761653436393937626134613333633937333434623730323137656266393034666538336166636230363062616432613562306365393433346564353961343536386430383234306462366412800135366266326535383864653163363264393463386631306263323538656263323433663038323735373336396537666363306664313331373461303533643039346138363362303532353662633065636438353665636633356236646462343935623431616433306565356333343837336666333036326433306331326331391abb030ab0020a423033643664386162393036613061643236333632386539633831663031613733646330333631633531633930613162353833616563653831303331323662663430631a0361626322033132332a0364656632033132333a03646566420a6d79656e636f64696e674a80013237383634636335323139613935316137613665353262386338646464663639383164303938646131363538643936323538633837306232633838646662636235313834316165613137326132386261666136613739373331313635353834363737303636303435633935396564306639393239363838643034646566633239524230336436643861623930366130616432363336323865396338316630316137336463303336316335316339306131623538336165636538313033313236626634306312800162313635303330326365333563306538626334323735376631346537303836336462636465393863313637346163613466376231303835316435383538636331363039653031623761303166393235363637366139396530646230356162313934353135346132326564656364663262363530386630326665643462336463621a030102031abb030ab0020a423033643664386162393036613061643236333632386539633831663031613733646330333631633531633930613162353833616563653831303331323662663430631a0361626322033132332a0364656632033435363a03676869420a6d79656e636f64696e674a80013237383634636335323139613935316137613665353262386338646464663639383164303938646131363538643936323538633837306232633838646662636235313834316165613137326132386261666136613739373331313635353834363737303636303435633935396564306639393239363838643034646566633239524230336436643861623930366130616432363336323865396338316630316137336463303336316335316339306131623538336165636538313033313236626634306312800162643630393032643031363239386334653263323830666538306137366362396634616261393666643235393363316632633566376165343639393762613461333363393733343462373032313765626639303466653833616663623036306261643261356230636539343334656435396134353638643038323430646236641a03010203"

    if datastr != expected {
        t.Error("Did not correctly encode batch. Got", datastr)
    }
}
