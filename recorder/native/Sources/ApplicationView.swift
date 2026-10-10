import SwiftUI

// Retained once by the generic application shell, including foreground recovery.
@MainActor final class ApplicationState: ObservableObject {
    let input: RecordingInput
    let controller: RecorderController?

    init(host: NativeHost) {
        let input = RecordingInput()
        self.input = input
        controller = host.url.map { RecorderController(baseURL: $0, input: input) }
        host.onShutdown = { [weak self] in self?.controller?.shutdown() }
        host.onBackgroundChange = { [weak self] backgrounded in
            self?.controller?.setAppBackgrounded(backgrounded)
        }
        host.onAvailabilityChange = { [weak self] available, message in
            // Only the WebKit receipt can confirm that transcript text is visible.
            if !available { self?.controller?.displayAvailability(false, error: message) }
        }
    }
}

struct ApplicationView: View {
    @ObservedObject var state: ApplicationState
    let url: URL
    let recoveryRevision: Int
    var body: some View {
        if let controller = state.controller {
            RecorderScreen(controller: controller, input: state.input, url: url, recoveryRevision: recoveryRevision)
        } else {
            Text("The recorder could not connect to the local application server.").padding()
        }
    }
}

#if os(macOS)
struct ApplicationSettingsView: View {
    @ObservedObject var state: ApplicationState
    var body: some View {
        RecordingSettings(input: state.input, controller: state.controller)
    }
}
#endif
