package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "strings"
    "time"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

type Node struct {
    Name  string
    X, Y  float64
    Edges []*Node
}

type Ant struct {
    ID            string
    Movements     []AntMovementStep
    PositionX     float64
    PositionY     float64
    FromNode      *Node
    ToNode        *Node
    Progress      float64 // Between 0 and 1
    Animating     bool
    CurrentNode   *Node
    MovementIndex int // Index in Movements
}

type AntMovementStep struct {
    Frame    int
    NodeName string
}

func main() {
    edges := readEdges("input.txt")
    graph := buildGraph(edges)
    assignPositions(graph, 400, 300, 200) // Center at (400,300), radius 200

    // Read ant movements
    antSequences, err := readAntMovements("movements.txt")
    if err != nil {
        fmt.Println(err)
        return
    }

    // Build ants map
    ants := make(map[string]*Ant)
    for antID := range antSequences {
        ant := &Ant{
            ID:            antID,
            Movements:     antSequences[antID],
            PositionX:     graph["start"].X,
            PositionY:     graph["start"].Y,
            Animating:     false,
            CurrentNode:   graph["start"],
            MovementIndex: 0,
        }
        ants[antID] = ant
    }

    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        fmt.Println("Error initializing SDL:", err)
        return
    }
    defer sdl.Quit()

    if err := ttf.Init(); err != nil {
        fmt.Println("Error initializing TTF:", err)
        return
    }
    defer ttf.Quit()

    window, err := sdl.CreateWindow("Graph Visualization", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 800, 600, sdl.WINDOW_SHOWN)
    if err != nil {
        fmt.Println("Error creating window:", err)
        return
    }
    defer window.Destroy()

    renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    if err != nil {
        fmt.Println("Error creating renderer:", err)
        return
    }
    defer renderer.Destroy()

    font, err := ttf.OpenFont("Arial.ttf", 16)
    if err != nil {
        fmt.Println("Failed to load font:", err)
        return
    }
    defer font.Close()

    const MovementDuration = 1000.0 // milliseconds per movement (animation duration)

    running := true
    advanceFrame := false
    animating := false
    frameIndex := 0
    lastUpdateTime := time.Now()

    // Prepare frames from movements.txt
    frames := parseFrames(antSequences)

    for running {
        // Handle events
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch e := event.(type) {
            case *sdl.QuitEvent:
                running = false
            case *sdl.KeyboardEvent:
                if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_RIGHT {
                    if !animating && frameIndex < len(frames) {
                        advanceFrame = true
                    }
                }
            }
        }

        currentTime := time.Now()
        deltaTime := currentTime.Sub(lastUpdateTime).Seconds() * 1000 // in milliseconds
        lastUpdateTime = currentTime

        if advanceFrame && frameIndex < len(frames) {
            // Start animating ants for the next frame
            frameMovements := frames[frameIndex]
            for antID, nodeName := range frameMovements {
                ant := ants[antID]
                if ant != nil {
                    ant.Animating = true
                    ant.Progress = 0.0
                    ant.FromNode = ant.CurrentNode
                    ant.ToNode = graph[nodeName]
                    ant.MovementIndex++
                }
            }
            animating = true
            advanceFrame = false
            frameIndex++
        }

        if animating {
            allAntsFinished := true
            for _, ant := range ants {
                if ant.Animating {
                    ant.Progress += deltaTime / MovementDuration
                    if ant.Progress >= 1.0 {
                        ant.PositionX = ant.ToNode.X
                        ant.PositionY = ant.ToNode.Y
                        ant.CurrentNode = ant.ToNode // Update the current node
                        ant.Animating = false
                    } else {
                        ant.PositionX = ant.FromNode.X + (ant.ToNode.X - ant.FromNode.X)*ant.Progress
                        ant.PositionY = ant.FromNode.Y + (ant.ToNode.Y - ant.FromNode.Y)*ant.Progress
                    }
                    allAntsFinished = false
                }
            }
            if allAntsFinished {
                animating = false
            }
        }

        // Update positions of ants that are not animating
        for _, ant := range ants {
            if !ant.Animating {
                ant.PositionX = ant.CurrentNode.X
                ant.PositionY = ant.CurrentNode.Y
            }
        }

        // Draw the graph with the current ant positions
        drawGraph(renderer, font, graph, ants)

        // Delay to control frame rate
        sdl.Delay(16) // Approximately 60 FPS
    }
}

func parseFrames(antSequences map[string][]AntMovementStep) []map[string]string {
    // Find the total number of frames
    maxFrame := 0
    for _, steps := range antSequences {
        for _, step := range steps {
            if step.Frame > maxFrame {
                maxFrame = step.Frame
            }
        }
    }

    // Initialize frames
    frames := make([]map[string]string, maxFrame+1)
    for i := 0; i <= maxFrame; i++ {
        frames[i] = make(map[string]string)
    }

    // Populate frames with ant movements
    for antID, steps := range antSequences {
        for _, step := range steps {
            frames[step.Frame][antID] = step.NodeName
        }
    }

    return frames
}

func readEdges(filename string) [][2]string {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return nil
    }
    defer file.Close()

    var edges [][2]string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        nodes := strings.Split(line, "-")
        if len(nodes) == 2 {
            edges = append(edges, [2]string{nodes[0], nodes[1]})
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
    }

    return edges
}

func readAntMovements(filename string) (map[string][]AntMovementStep, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("Error opening ant movements file: %v", err)
    }
    defer file.Close()

    antSequences := make(map[string][]AntMovementStep)
    scanner := bufio.NewScanner(file)
    frame := 0

    for scanner.Scan() {
        line := scanner.Text()
        tokens := strings.Fields(line)
        for _, token := range tokens {
            parts := strings.Split(token, "-")
            if len(parts) == 2 {
                antID := parts[0]
                nodeName := parts[1]
                antSequences[antID] = append(antSequences[antID], AntMovementStep{
                    Frame:    frame,
                    NodeName: nodeName,
                })
            }
        }
        frame++
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("Error reading ant movements file: %v", err)
    }

    return antSequences, nil
}

func buildGraph(edges [][2]string) map[string]*Node {
    nodes := make(map[string]*Node)

    for _, edge := range edges {
        fromName, toName := edge[0], edge[1]

        fromNode, exists := nodes[fromName]
        if !exists {
            fromNode = &Node{Name: fromName}
            nodes[fromName] = fromNode
        }

        toNode, exists := nodes[toName]
        if !exists {
            toNode = &Node{Name: toName}
            nodes[toName] = toNode
        }

        // Directed edge from 'fromNode' to 'toNode'
        fromNode.Edges = append(fromNode.Edges, toNode)
    }

    return nodes
}

func assignPositions(nodes map[string]*Node, centerX, centerY, radius float64) {
    totalNodes := len(nodes)
    angleIncrement := 2 * math.Pi / float64(totalNodes)
    i := 0

    for _, node := range nodes {
        angle := angleIncrement * float64(i)
        node.X = centerX + radius*math.Cos(angle)
        node.Y = centerY + radius*math.Sin(angle)
        i++
    }
}

func drawGraph(renderer *sdl.Renderer, font *ttf.Font, nodes map[string]*Node, ants map[string]*Ant) {
    renderer.SetDrawColor(255, 255, 255, 255) // White background
    renderer.Clear()

    // Draw edges
    for _, node := range nodes {
        for _, neighbor := range node.Edges {
            drawArrow(renderer, node.X, node.Y, neighbor.X, neighbor.Y)
        }
    }

    // Draw nodes and labels
    for _, node := range nodes {
        drawCircle(renderer, int32(node.X), int32(node.Y), 20)
        drawLabel(renderer, font, node)
    }

    // Draw ants
    for _, ant := range ants {
        drawAnt(renderer, int32(ant.PositionX), int32(ant.PositionY))
    }

    renderer.Present()
}

func drawCircle(renderer *sdl.Renderer, x0, y0, radius int32) {
    // Draw filled circle (light yellow)
    renderer.SetDrawColor(255, 255, 224, 255) // Light yellow fill color
    for w := -radius; w <= radius; w++ {
        for h := -radius; h <= radius; h++ {
            if w*w+h*h <= radius*radius {
                renderer.DrawPoint(x0+w, y0+h)
            }
        }
    }

    // Draw circle outline (black)
    renderer.SetDrawColor(0, 0, 0, 255) // Black outline color
    for angle := 0.0; angle <= 2*math.Pi; angle += 0.01 {
        x := x0 + int32(float64(radius)*math.Cos(angle))
        y := y0 + int32(float64(radius)*math.Sin(angle))
        renderer.DrawPoint(x, y)
    }
}

func drawArrow(renderer *sdl.Renderer, x1, y1, x2, y2 float64) {
    renderer.SetDrawColor(0, 0, 0, 255) // Black color for edges

    // Draw the line
    renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

    // Calculate the angle of the line
    angle := math.Atan2(y2 - y1, x2 - x1)

    // Arrowhead size
    arrowSize := 10.0

    // Calculate the points for the arrowhead
    x3 := x2 - arrowSize*math.Cos(angle - math.Pi/6)
    y3 := y2 - arrowSize*math.Sin(angle - math.Pi/6)
    x4 := x2 - arrowSize*math.Cos(angle + math.Pi/6)
    y4 := y2 - arrowSize*math.Sin(angle + math.Pi/6)

    // Draw the arrowhead
    renderer.DrawLine(int32(x2), int32(y2), int32(x3), int32(y3))
    renderer.DrawLine(int32(x2), int32(y2), int32(x4), int32(y4))
}

func drawLabel(renderer *sdl.Renderer, font *ttf.Font, node *Node) {
    // Render the text to a surface
    surface, err := font.RenderUTF8Solid(node.Name, sdl.Color{R: 0, G: 0, B: 0, A: 255})
    if err != nil {
        fmt.Println("Failed to render text:", err)
        return
    }
    defer surface.Free()

    // Create a texture from the surface
    texture, err := renderer.CreateTextureFromSurface(surface)
    if err != nil {
        fmt.Println("Failed to create texture:", err)
        return
    }
    defer texture.Destroy()

    // Get the dimensions of the text
    textWidth := surface.W
    textHeight := surface.H

    // Define the rectangle where the text will be drawn
    dstRect := sdl.Rect{
        X: int32(node.X) - textWidth/2,
        Y: int32(node.Y) - textHeight/2,
        W: textWidth,
        H: textHeight,
    }

    // Copy the texture to the renderer
    renderer.Copy(texture, nil, &dstRect)
}

func drawAnt(renderer *sdl.Renderer, x, y int32) {
    renderer.SetDrawColor(255, 0, 0, 255) // Red color for ants
    radius := int32(5)                    // Ant size

    for w := -radius; w <= radius; w++ {
        for h := -radius; h <= radius; h++ {
            if w*w+h*h <= radius*radius {
                renderer.DrawPoint(x+w, y+h)
            }
        }
    }
}
