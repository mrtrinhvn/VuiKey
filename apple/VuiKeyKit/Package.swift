// swift-tools-version: 5.7
import PackageDescription

let package = Package(
    name: "VuiKeyKit",
    platforms: [
        .macOS(.v12),
        .iOS(.v15)
    ],
    products: [
        .library(
            name: "VuiKeyEngine",
            targets: ["VuiKeyEngine"]
        )
    ],
    targets: [
        .target(
            name: "VuiKeyEngine",
            dependencies: []
        ),
        .testTarget(
            name: "VuiKeyEngineTests",
            dependencies: ["VuiKeyEngine"]
        )
    ]
)
