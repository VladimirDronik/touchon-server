package ModbusGW

type regType int

const (
	Coil regType = iota
	Hold
)

type register struct {
	Type    regType
	Address uint16
}

type Gateway struct {
	Name        string
	OpModes     map[string]string
	FanSpeed    map[string]string
	HSlatsModes map[string]string
	VSlatsModes map[string]string
	Registers   map[string]*register
}

var gateways = map[string]*Gateway{
	// https://onokom.ru/GR-1-MB-B.html
	"gr_1_mb_b": {
		Name: "GR-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Левое положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     {Coil, 0x0006},
			"turbo_mode":                   {Coil, 0x0007},
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   {Coil, 0x0009},
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": {Coil, 0x0029},
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/GR-3-MB-B.html
	"gr_3_mb_b": {
		Name: "GR-3-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
			"4": "Четвертая скорость",
			"5": "Пятая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Левое положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            nil,
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     nil,
			"turbo_mode":                   nil,
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   nil,
			"self_cleaning":                {Coil, 0x000A},
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/TCL-1-MB-B.html
	"tcl_1_mb_b": {
		Name: "TCL-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
			"4": "Четвертая скорость",
			"5": "Пятая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Левое положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
			"7": "Мягкий поток",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     {Coil, 0x0006},
			"turbo_mode":                   {Coil, 0x0007},
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   {Coil, 0x0009},
			"self_cleaning":                {Coil, 0x000A},
			"anti_fungus":                  {Coil, 0x000B},
			"disable_display_on_power_off": nil,
			"sounds":                       {Coil, 0x0005},
			"on_duty_heating":              {Coil, 0x000D},
			"soft_flow":                    {Coil, 0x000E},
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/DK-1-MB-B.html
	"dk_1_mb_b": {
		Name: "DK-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Тихий режим",
			"2": "Первая скорость",
			"3": "Вторая скорость",
			"4": "Третья скорость",
			"5": "Четвертая скорость",
			"6": "Пятая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      {Coil, 0x0003},
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     {Coil, 0x0006},
			"turbo_mode":                   {Coil, 0x0007},
			"sleep_mode":                   nil,
			"ionization":                   nil,
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/AUX-1-MB-B.html
	"aux_1_mb_b": {
		Name: "AUX-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Тихий режим",
			"2": "Первая скорость",
			"3": "Вторая скорость",
			"4": "Третья скорость",
			"5": "Четвертая скорость",
			"6": "Пятая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     {Coil, 0x0006},
			"turbo_mode":                   {Coil, 0x0007},
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   {Coil, 0x0009},
			"self_cleaning":                {Coil, 0x000A},
			"anti_fungus":                  {Coil, 0x000B},
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/ME-1-MB-B.html#karta-registrov
	"me_1_mb_b": {
		Name: "ME-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
			"4": "Четвертая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Левое положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            nil,
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     nil,
			"turbo_mode":                   nil,
			"sleep_mode":                   nil,
			"ionization":                   nil,
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         nil,
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/HS-3-MB-B.html#karta-registrov
	"hs_3_mb_b": {
		Name: "HS-3-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
			"4": "Четвертая скорость",
			"5": "Пятая скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0004},
			"eco_mode":                     {Coil, 0x0006},
			"turbo_mode":                   nil,
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   nil,
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/HR-1-MB-B.html
	"hr_1_mb_b": {
		Name: "HR-1-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
		},
		HSlatsModes: map[string]string{
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
			"7": "Шестое положение",
			"8": "Седьмое положение",
		},
		VSlatsModes: map[string]string{
			"1": "Качание",
			"2": "Левое положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
			"6": "Пятое положение",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            {Coil, 0x0002},
			"display_high_brightness":      nil,
			"silent_mode":                  {Coil, 0x0003},
			"eco_mode":                     {Coil, 0x0005},
			"turbo_mode":                   {Coil, 0x0006},
			"sleep_mode":                   {Coil, 0x0007},
			"ionization":                   {Coil, 0x0008},
			"self_cleaning":                {Coil, 0x000A},
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       {Coil, 0x0004},
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         {Hold, 0x0004},
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/HS-6-MB-B.html#karta-registrov
	"hs_6_mb_b": {
		Name: "HS-6-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		VSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
		},
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            nil,
			"display_high_brightness":      nil,
			"silent_mode":                  nil,
			"eco_mode":                     nil,
			"turbo_mode":                   nil,
			"sleep_mode":                   {Coil, 0x0008},
			"ionization":                   nil,
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         nil,
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          {Hold, 0x000A},
		},
	},

	// https://onokom.ru/MH-8-MB-B.html
	"mh_8_mb_b": {
		Name: "MH-8-MB-B",
		OpModes: map[string]string{
			"1": "Нагрев",
			"2": "Охлаждение",
			"3": "Автоматический",
			"4": "Осушение",
			"5": "Вентиляция",
		},
		FanSpeed: map[string]string{
			"0": "Авто",
			"1": "Первая скорость",
			"2": "Вторая скорость",
			"3": "Третья скорость",
		},
		HSlatsModes: map[string]string{
			"0": "Остановлено",
			"1": "Качание",
			"2": "Нижнее положение",
			"3": "Второе положение",
			"4": "Третье положение",
			"5": "Четвертое положение",
		},
		VSlatsModes: nil,
		Registers: map[string]*register{
			"power_status":                 {Coil, 0x0001},
			"display_backlight":            nil,
			"display_high_brightness":      nil,
			"silent_mode":                  nil,
			"eco_mode":                     nil,
			"turbo_mode":                   nil,
			"sleep_mode":                   nil,
			"ionization":                   nil,
			"self_cleaning":                nil,
			"anti_fungus":                  nil,
			"disable_display_on_power_off": nil,
			"sounds":                       nil,
			"on_duty_heating":              nil,
			"soft_flow":                    nil,
			"operating_mode":               {Hold, 0x0001},
			"internal_temperature":         {Hold, 0x0003},
			"external_temperature":         nil,
			"target_temperature":           {Hold, 0x0005},
			"fan_speed":                    {Hold, 0x0007},
			"horizontal_slats_mode":        {Hold, 0x0009},
			"vertical_slats_mode":          nil,
		},
	},
}

func (o *Gateway) getMaxCoilAddress() uint16 {
	res := uint16(0)

	for _, r := range o.Registers {
		if r != nil && r.Type == Coil && r.Address > res {
			res = r.Address
		}
	}

	return res
}

func (o *Gateway) getMaxHoldAddress() uint16 {
	res := uint16(0)

	for _, r := range o.Registers {
		if r != nil && r.Type == Hold && r.Address > res {
			res = r.Address
		}
	}

	return res
}

// Экспортируем список поддерживаемых шлюзов

var SupportedGateways = map[string]string{}

func init() {
	for gwModelCode, gw := range gateways {
		SupportedGateways[gwModelCode] = gw.Name
	}
}
