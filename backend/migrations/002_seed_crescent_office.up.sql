-- Clear existing seed data before inserting updated records
TRUNCATE TABLE user_progress, hotspots, quests, zones CASCADE;

-- ==========================================
-- 1. SEED ZONES (Strictly matching latest Xsolla.svg IDs)
-- ==========================================
INSERT INTO zones (id, name, is_locked_default, difficulty, media_url, description) VALUES
    -- Game Room Sectors (1 = Easy -> 4 = Hard)
    (
        'zone_gameroom_1',
        'Game Room — Sector 1 (Board Games)',
        TRUE,
        'EASY',
        '/assets/rooms/gameroom_sector_1.jpg',
        'Board games and tabletop lounge area.'
    ),
    (
        'zone_gameroom_2',
        'Game Room — Sector 2 (Foosball & Table Hockey)',
        TRUE,
        'MEDIUM',
        '/assets/rooms/gameroom_sector_2.jpg',
        'Foosball and table hockey zone.'
    ),
    (
        'zone_gameroom_3',
        'Game Room — Sector 3 (Table Tennis)',
        TRUE,
        'MEDIUM_HARD',
        '/assets/rooms/gameroom_sector_3.jpg',
        'Table tennis zone.'
    ),
    (
        'zone_gameroom_4',
        'Game Room — Sector 4 (PlayStation Zone)',
        TRUE,
        'HARD',
        '/assets/rooms/gameroom_sector_4.jpg',
        'PlayStation console gaming zone.'
    ),

    -- 11 Meeting Rooms (zone_meeting_1 .. zone_meeting_11)
    (
        'zone_meeting_1',
        'Meeting Room 1',
        TRUE,
        'EASY',
        '/assets/rooms/meeting_1.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_2',
        'Meeting Room 2',
        TRUE,
        'EASY',
        '/assets/rooms/meeting_2.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_3',
        'Meeting Room 3',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_3.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_4',
        'Meeting Room 4',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_4.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_5',
        'Meeting Room 5',
        TRUE,
        'EASY',
        '/assets/rooms/meeting_5.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_6',
        'Meeting Room 6',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_6.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_7',
        'Meeting Room 7',
        TRUE,
        'EASY',
        '/assets/rooms/meeting_7.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_8',
        'Meeting Room 8',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_8.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),
    (
        'zone_meeting_9',
        'Meeting Room 9',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_9.jpg',
        'Presentation meeting room equipped with a projector. Book via calendar before use.'
    ),
    (
        'zone_meeting_10',
        'Meeting Room 10',
        TRUE,
        'MEDIUM',
        '/assets/rooms/meeting_10.jpg',
        'Presentation meeting room equipped with a projector. Book via calendar before use.'
    ),
    (
        'zone_meeting_11',
        'Meeting Room 11',
        TRUE,
        'EASY',
        '/assets/rooms/meeting_11.jpg',
        'Meeting room for team syncs and calls. Book via calendar before use.'
    ),

    -- Unlocked Zones (Open Spaces, Session Hall, Kitchen, Stock)
    (
        'zone_lecturezone_open_space_1',
        'Open Space 1 (Xsolla School Students)',
        FALSE,
        NULL,
        '/assets/rooms/open_space_1.jpg',
        'Dedicated open workspace for Xsolla School students.'
    ),
    (
        'zone_open_space_2',
        'Open Space 2 (Xsolla Employees)',
        FALSE,
        NULL,
        '/assets/rooms/open_space_2.jpg',
        'Main open workspace area for Xsolla employees.'
    ),
    (
        'zone_session',
        'Session & Lecture Hall',
        FALSE,
        NULL,
        '/assets/rooms/session_hall.jpg',
        'Session and lecture hall for presentations, workshops, and talks.'
    ),
    (
        'zone_kitchen',
        'Kitchen',
        FALSE,
        NULL,
        '/assets/rooms/kitchen.jpg',
        'Office kitchen equipped with a coffee machine, refrigerator, and water dispenser.'
    ),
    (
        'zone_stock',
        'Stock Room',
        FALSE,
        NULL,
        '/assets/rooms/stock.jpg',
        'Office storage and equipment supply room.'
    );

-- ==========================================
-- 2. SEED ALGORITHMIC QUESTS (1:1 with Locked Zones)
-- ==========================================
INSERT INTO quests (zone_id, title, problem_statement, code_blocks, correct_order) VALUES
    -- Game Room 1 (EASY - Board Games)
    (
        'zone_gameroom_1',
        'Sector 1 (Easy): Swap Two Values',
        'Arrange the Go code blocks in the correct order to swap two integer values a and b.',
        '[
            {"id": "b1", "code": "func Swap(a, b int) (int, int) {"},
            {"id": "b2", "code": "    temp := a"},
            {"id": "b3", "code": "    a = b"},
            {"id": "b4", "code": "    b = temp"},
            {"id": "b5", "code": "    return a, b\n}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    -- Game Room 2 (MEDIUM - Foosball & Table Hockey)
    (
        'zone_gameroom_2',
        'Sector 2 (Medium): Find Maximum Score',
        'Arrange the code blocks to find and return the maximum value in an integer slice.',
        '[
            {"id": "b1", "code": "func MaxScore(scores []int) int {"},
            {"id": "b2", "code": "    maxVal := scores[0]"},
            {"id": "b3", "code": "    for _, s := range scores[1:] {"},
            {"id": "b4", "code": "        if s > maxVal { maxVal = s }"},
            {"id": "b5", "code": "    }"},
            {"id": "b6", "code": "    return maxVal\n}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5", "b6"]'::jsonb
    ),
    -- Game Room 3 (MEDIUM_HARD - Table Tennis)
    (
        'zone_gameroom_3',
        'Sector 3 (Medium-Hard): Binary Search',
        'Order the blocks to perform a binary search on a sorted slice of integers.',
        '[
            {"id": "b1", "code": "low, high := 0, len(arr)-1"},
            {"id": "b2", "code": "for low <= high {"},
            {"id": "b3", "code": "    mid := low + (high-low)/2"},
            {"id": "b4", "code": "    if arr[mid] == target { return mid }"},
            {"id": "b5", "code": "    if arr[mid] < target { low = mid + 1 } else { high = mid - 1 }"},
            {"id": "b6", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5", "b6"]'::jsonb
    ),
    -- Game Room 4 (HARD - PlayStation Zone)
    (
        'zone_gameroom_4',
        'Sector 4 (Hard): Worker Pool Dispatcher',
        'Arrange the Go concurrency blocks to start workers, send jobs, close the channel, and wait for completion.',
        '[
            {"id": "b1", "code": "jobs := make(chan int, len(tasks))"},
            {"id": "b2", "code": "var wg sync.WaitGroup"},
            {"id": "b3", "code": "for w := 0; w < numWorkers; w++ {\n    wg.Add(1)\n    go worker(jobs, &wg)\n}"},
            {"id": "b4", "code": "for _, t := range tasks {\n    jobs <- t\n}"},
            {"id": "b5", "code": "close(jobs)"},
            {"id": "b6", "code": "wg.Wait()"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5", "b6"]'::jsonb
    ),

    -- Meeting Rooms 1..11 Quests
    (
        'zone_meeting_1',
        'Meeting Room 1: JSON Response Helper',
        'Order the statements to set the JSON header, write the status code, and encode the response.',
        '[
            {"id": "b1", "code": "func WriteJSON(w http.ResponseWriter, status int, data any) error {"},
            {"id": "b2", "code": "    w.Header().Set(\"Content-Type\", \"application/json\")"},
            {"id": "b3", "code": "    w.WriteHeader(status)"},
            {"id": "b4", "code": "    return json.NewEncoder(w).Encode(data)"},
            {"id": "b5", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_2',
        'Meeting Room 2: Time Slot Conflict Check',
        'Arrange the function to check if two meeting intervals overlap.',
        '[
            {"id": "b1", "code": "func HasConflict(startA, endA, startB, endB int) bool {"},
            {"id": "b2", "code": "    if endA <= startB || endB <= startA {"},
            {"id": "b3", "code": "        return false"},
            {"id": "b4", "code": "    }"},
            {"id": "b5", "code": "    return true\n}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_3',
        'Meeting Room 3: Stack Push and Pop',
        'Order the lines to push an item onto a slice stack and pop the last item.',
        '[
            {"id": "b1", "code": "stack := []int{10, 20}"},
            {"id": "b2", "code": "stack = append(stack, 30)"},
            {"id": "b3", "code": "top := stack[len(stack)-1]"},
            {"id": "b4", "code": "stack = stack[:len(stack)-1]"},
            {"id": "b5", "code": "fmt.Println(top)"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_4',
        'Meeting Room 4: Context Timeout',
        'Arrange the code to create a timeout context, defer cancel, and ping the database.',
        '[
            {"id": "b1", "code": "func PingDB(ctx context.Context, db *sql.DB) error {"},
            {"id": "b2", "code": "    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)"},
            {"id": "b3", "code": "    defer cancel()"},
            {"id": "b4", "code": "    return db.PingContext(ctx)"},
            {"id": "b5", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_5',
        'Meeting Room 5: Word Counter',
        'Order the blocks to count word occurrences in a slice.',
        '[
            {"id": "b1", "code": "counts := make(map[string]int)"},
            {"id": "b2", "code": "for _, w := range words {"},
            {"id": "b3", "code": "    counts[w]++"},
            {"id": "b4", "code": "}"},
            {"id": "b5", "code": "return counts"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_6',
        'Meeting Room 6: Mutex Counter',
        'Arrange the statements to safely increment a shared counter with a mutex.',
        '[
            {"id": "b1", "code": "func (c *Counter) Inc() {"},
            {"id": "b2", "code": "    c.mu.Lock()"},
            {"id": "b3", "code": "    defer c.mu.Unlock()"},
            {"id": "b4", "code": "    c.value++"},
            {"id": "b5", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_7',
        'Meeting Room 7: Filter Even Numbers',
        'Order the blocks to filter even numbers from a slice.',
        '[
            {"id": "b1", "code": "var evens []int"},
            {"id": "b2", "code": "for _, n := range nums {"},
            {"id": "b3", "code": "    if n%2 == 0 { evens = append(evens, n) }"},
            {"id": "b4", "code": "}"},
            {"id": "b5", "code": "return evens"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_8',
        'Meeting Room 8: Palindrome Check',
        'Arrange the blocks to check if a string is a palindrome using two pointers.',
        '[
            {"id": "b1", "code": "left, right := 0, len(s)-1"},
            {"id": "b2", "code": "for left < right {"},
            {"id": "b3", "code": "    if s[left] != s[right] { return false }"},
            {"id": "b4", "code": "    left++; right--"},
            {"id": "b5", "code": "}"},
            {"id": "b6", "code": "return true"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5", "b6"]'::jsonb
    ),
    (
        'zone_meeting_9',
        'Meeting Room 9: Reverse a Slice',
        'Arrange the blocks to reverse an integer slice in place.',
        '[
            {"id": "b1", "code": "func Reverse(s []int) {"},
            {"id": "b2", "code": "    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {"},
            {"id": "b3", "code": "        s[i], s[j] = s[j], s[i]"},
            {"id": "b4", "code": "    }"},
            {"id": "b5", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_10',
        'Meeting Room 10: Sum of Slice',
        'Arrange the blocks to calculate and return the sum of all integers in a slice.',
        '[
            {"id": "b1", "code": "func Sum(nums []int) int {"},
            {"id": "b2", "code": "    total := 0"},
            {"id": "b3", "code": "    for _, n := range nums { total += n }"},
            {"id": "b4", "code": "    return total"},
            {"id": "b5", "code": "}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    ),
    (
        'zone_meeting_11',
        'Meeting Room 11: Nil Slice Check',
        'Order the statements to return an error if the input slice is empty, or return the first element.',
        '[
            {"id": "b1", "code": "func First(items []string) (string, error) {"},
            {"id": "b2", "code": "    if len(items) == 0 {"},
            {"id": "b3", "code": "        return \"\", errors.New(\"empty slice\")"},
            {"id": "b4", "code": "    }"},
            {"id": "b5", "code": "    return items[0], nil\n}"}
        ]'::jsonb,
        '["b1", "b2", "b3", "b4", "b5"]'::jsonb
    );

-- ==========================================
-- 3. SEED HOTSPOTS & EASTER EGGS (Matching latest Xsolla.svg)
-- ==========================================
INSERT INTO hotspots (id, zone_id, category, x_percent, y_percent, title, content) VALUES
    -- 3 Water Dispensers from Xsolla.svg
    (
        'hotspot_dispenser_1',
        'zone_lecturezone_open_space_1',
        'DISPENSER',
        42.34,
        38.05,
        'Water Dispenser (School Open Space)',
        'Filtered cold and ambient water dispenser.'
    ),
    (
        'hotspot_dispenser_2',
        'zone_lecturezone_open_space_1',
        'DISPENSER',
        62.08,
        37.78,
        'Water Dispenser (Corridor)',
        'Filtered cold and ambient water dispenser.'
    ),
    (
        'hotspot_dispenser_3',
        'zone_kitchen',
        'DISPENSER',
        68.12,
        53.70,
        'Kitchen Water Dispenser',
        'Filtered cold, ambient, and hot water dispenser.'
    ),

    -- Kitchen Equipment
    (
        'hotspot_kitchen_coffee',
        'zone_kitchen',
        'COFFEE_MACHINE',
        32.00,
        52.00,
        'Coffee Machine',
        'Place your cup under the spout and select your coffee on the touch screen.'
    ),
    (
        'hotspot_kitchen_fridge',
        'zone_kitchen',
        'FRIDGE',
        64.00,
        48.00,
        'Refrigerator',
        'Store labeled food containers here. Unlabeled items are cleared on Fridays.'
    ),

    -- 3 Projectors from Xsolla.svg (zone_session, zone_meeting_9, zone_meeting_10)
    (
        'hotspot_projector_1',
        'zone_session',
        'PROJECTOR',
        65.28,
        42.13,
        'Session Hall Projector',
        'Connect via HDMI or wireless screen sharing for presentations.'
    ),
    (
        'hotspot_projector_2',
        'zone_meeting_9',
        'PROJECTOR',
        49.66,
        67.40,
        'Meeting Room 9 Projector',
        'Connect via HDMI or wireless screen sharing.'
    ),
    (
        'hotspot_projector_3',
        'zone_meeting_10',
        'PROJECTOR',
        45.34,
        67.78,
        'Meeting Room 10 Projector',
        'Connect via HDMI or wireless screen sharing.'
    ),

    -- Viewpoint
    (
        'hotspot_viewpoint',
        NULL,
        'VIEWPOINT',
        88.00,
        45.00,
        'Crescent Bay Viewpoint',
        'Panoramic view of Crescent Bay and the Baku skyline from the 31st floor.'
    ),

    -- Hidden Easter Eggs
    (
        'hotspot_egg_gameroom_1',
        'zone_gameroom_1',
        'EASTER_EGG',
        76.50,
        68.20,
        'Hidden Xsolla Logo — Game Room 1',
        'You found the hidden Xsolla logo in the Board Games zone!'
    ),
    (
        'hotspot_egg_gameroom_3',
        'zone_gameroom_3',
        'EASTER_EGG',
        22.10,
        33.80,
        'Hidden Xsolla Logo — Game Room 3',
        'You found the hidden Xsolla logo in the Table Tennis zone!'
    ),
    (
        'hotspot_egg_kitchen',
        'zone_kitchen',
        'EASTER_EGG',
        82.30,
        24.50,
        'Hidden Xsolla Logo — Kitchen',
        'You found the hidden Xsolla logo in the Kitchen!'
    ),
    (
        'hotspot_egg_session',
        'zone_session',
        'EASTER_EGG',
        14.20,
        78.40,
        'Hidden Xsolla Logo — Session Hall',
        'You found the hidden Xsolla logo in the Session Hall!'
    ),
    (
        'hotspot_egg_viewpoint',
        NULL,
        'EASTER_EGG',
        91.40,
        48.20,
        'Hidden Xsolla Logo — Viewpoint',
        'You found the hidden Xsolla logo at the Viewpoint!'
    );
