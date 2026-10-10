import Foundation
import AppleTranscription

@main struct RecordingInputTests {
    @MainActor static func main() {
        var inventoryCalls = 0
        let inventory: @MainActor () throws -> [InputDevice] = {
            inventoryCalls += 1
            return [InputDevice(uid: "test-input", name: "Test microphone")]
        }
        for value in ["fake", "synthetci", "", "REAL"] {
            let invalid = RecordingInput(environment: ["VN_SOURCE": value], listInputs: inventory)
            precondition(invalid.configurationError != nil && invalid.error == invalid.configurationError)
            precondition(invalid.fixtureURL == nil && invalid.devices.isEmpty)
            invalid.refresh()
            precondition(inventoryCalls == 0 && invalid.error == invalid.configurationError)
        }
        let scripted = RecordingInput(environment: ["VN_SOURCE": "synthetic"], listInputs: inventory)
        precondition(scripted.fixtureURL != nil && scripted.configurationError == nil)
        scripted.refresh()
        precondition(inventoryCalls == 0)
        // An explicit fixture wins even when VN_SOURCE is invalid.
        let fixture = RecordingInput(environment: ["VN_SOURCE": "fake", "VN_TRANSCRIPT_FIXTURE": "/tmp/explicit-fixture.json"], listInputs: inventory)
        precondition(fixture.fixtureURL?.path == "/tmp/explicit-fixture.json" && fixture.configurationError == nil)
        precondition(inventoryCalls == 0)
        let real = RecordingInput(environment: ["VN_SOURCE": "real", "VN_INPUT_DEVICE": "test-input"], listInputs: inventory)
        precondition(real.fixtureURL == nil && real.configurationError == nil && real.error == nil)
        precondition(real.selection == .device(uid: "test-input") && real.devices.count == 1 && inventoryCalls == 1)
        real.refresh()
        precondition(inventoryCalls == 2)
        print("PASS invalid sources never inventory microphones, synthetic and explicit fixtures avoid inventory, real selection retains pinned UID")
    }
}
