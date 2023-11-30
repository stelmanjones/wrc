package wrc

import (
	"time"

	"github.com/wI2L/jettison"
)

// Percentage is an alias for int
type Percentage int

// Packet represents the default WRC data packet.
type Packet struct {
	// A rolling unique identifier for the current packet. Can be used to order and drop received packets.
	PacketUID uint64 `json:"packet_uid"`

	// Time spent in game since boot.
	GameTotalTime float32 `json:"game_total_time"`

	// Time spent since last frame.
	GameDeltaTime float32 `json:"game_delta_time"`

	// Frame count in game since boot.
	GameFrameCount uint64 `json:"game_frame_count"`

	// For shift lights, from 0 ('vehicle_engine_rpm_current'='shiftlights_rpm_start') to 1 ('vehicle_engine_rpm_current'='shiftlights_rpm_end').
	ShiftlightsFraction float32 `json:"shiftlights_fraction"`

	// Shift lights start at 'vehicle_engine_rpm_current' value.
	ShiftlightsRpmStart float32 `json:"shiftlights_rpm_start"`

	// Shift lights end (i.e. optimal shift) at 'vehicle_engine_rpm_current' value.
	ShiftlightsRpmEnd float32 `json:"shiftlights_rpm_end"`

	// Are shift lights RPM data valid: 'vehicle_engine_rpm_current', 'shiftlights_rpm_start', 'shiftlights_rpm_end'
	ShiftlightsRpmValid bool `json:"shiftlights_rpm_valid"`

	// Gear index or value of 'vehicle_gear_index_neutral' or 'vehicle_gear_index_reverse'
	VehicleGearIndex uint8 `json:"vehicle_gear_index"`

	// 'vehicle_gear_index' if gearbox in Neutral.
	VehicleGearIndexNeutral uint8 `json:"vehicle_gear_index_neutral"`

	// 'vehicle_gear_index' if gearbox in Reverse.
	VehicleGearIndexReverse uint8 `json:"vehicle_gear_index_reverse"`

	// Number of forward gears.
	VehicleGearMaximum uint8 `json:"vehicle_gear_maximum"`

	// Car body speed.
	VehicleSpeed float32 `json:"vehicle_speed"`

	// Car speed at wheel/road due to transmission (for speedo use). NB. May differ from 'vehicle_speed'.
	VehicleTransmissionSpeed float32 `json:"vehicle_transmission_speed"`

	// Car position X component, positive left.
	VehiclePositionX float32 `json:"vehicle_position_x"`

	// Car position Y component, positive up.
	VehiclePositionY float32 `json:"vehicle_position_y"`

	// Car position Z component, positive forward.
	VehiclePositionZ float32 `json:"vehicle_position_z"`

	// Car velocity X component, positive left.
	VehicleVelocityX float32 `json:"vehicle_velocity_x"`

	// Car velocity Y component, positive up.
	VehicleVelocityY float32 `json:"vehicle_velocity_y"`

	// Car velocity Z component, positive forward.
	VehicleVelocityZ float32 `json:"vehicle_velocity_z"`

	// Car acceleration X component, positive left.
	VehicleAccelerationX float32 `json:"vehicle_acceleration_x"`

	// Car acceleration Y component, positive up.
	VehicleAccelerationY float32 `json:"vehicle_acceleration_y"`

	// Car acceleration Z component, positive forward.
	VehicleAccelerationZ float32 `json:"vehicle_acceleration_z"`

	// Car left unit vector X component, positive left.
	VehicleLeftDirectionX float32 `json:"vehicle_left_direction_x"`

	// Car left unit vector Y component, positive up.
	VehicleLeftDirectionY float32 `json:"vehicle_left_direction_y"`

	// Car left unit vector Z component, positive forward.
	VehicleLeftDirectionZ float32 `json:"vehicle_left_direction_z"`

	// Car forward unit vector X component, positive left.
	VehicleForwardDirectionX float32 `json:"vehicle_forward_direction_x"`

	// Car forward unit vector Y component, positive up.
	VehicleForwardDirectionY float32 `json:"vehicle_forward_direction_y"`

	// Car forward unit vector Z component, positive forward.
	VehicleForwardDirectionZ float32 `json:"vehicle_forward_direction_z"`

	// Car up unit vector X component, positive left.
	VehicleUpDirectionX float32 `json:"vehicle_up_direction_x"`

	// Car up unit vector Y component, positive up.
	VehicleUpDirectionY float32 `json:"vehicle_up_direction_y"`

	// Car up unit vector Z component, positive forward.
	VehicleUpDirectionZ float32 `json:"vehicle_up_direction_z"`

	// Wheel hub height displacement, back left, positive up.
	VehicleHubPositionBL float32 `json:"vehicle_hub_position_bl"`

	// Wheel hub height displacement, back right, positive up.
	VehicleHubPositionBR float32 `json:"vehicle_hub_position_br"`

	// Wheel hub height displacement, front left, positive up.
	VehicleHubPositionFL float32 `json:"vehicle_hub_position_fl"`

	// Wheel hub height displacement, front right, positive up.
	VehicleHubPositionFR float32 `json:"vehicle_hub_position_fr"`

	// Wheel hub vertical velocity, back left, positive up.
	VehicleHubVelocityBL float32 `json:"vehicle_hub_velocity_bl"`

	// Wheel hub vertical velocity, back right, positive up.
	VehicleHubVelocityBR float32 `json:"vehicle_hub_velocity_br"`

	// Wheel hub vertical velocity, front left, positive up.
	VehicleHubVelocityFL float32 `json:"vehicle_hub_velocity_fl"`

	// Wheel hub vertical velocity, front right, positive up.
	VehicleHubVelocityFR float32 `json:"vehicle_hub_velocity_fr"`

	// Contact patch forward speed, back left.
	VehicleCpForwardSpeedBL float32 `json:"vehicle_cp_forward_speed_bl"`

	// Contact patch forward speed, back right.
	VehicleCpForwardSpeedBR float32 `json:"vehicle_cp_forward_speed_br"`

	// Contact patch forward speed, front left.
	VehicleCpForwardSpeedFL float32 `json:"vehicle_cp_forward_speed_fl"`

	// Contact patch forward speed, front right.
	VehicleCpForwardSpeedFR float32 `json:"vehicle_cp_forward_speed_fr"`

	// Brake temperature, back left.
	VehicleBrakeTemperatureBL float32 `json:"vehicle_brake_temperature_bl"`

	// Brake temperature, back right.
	VehicleBrakeTemperatureBR float32 `json:"vehicle_brake_temperature_br"`

	// Brake temperature, front left.
	VehicleBrakeTemperatureFL float32 `json:"vehicle_brake_temperature_fl"`

	// Brake temperature, front right.
	VehicleBrakeTemperatureFR float32 `json:"vehicle_brake_temperature_fr"`

	// Engine rotation rate, maximum.
	VehicleEngineRpmMax float32 `json:"vehicle_engine_rpm_max"`

	// Engine rotation rate, at idle.
	VehicleEngineRpmIdle float32 `json:"vehicle_engine_rpm_idle"`

	// Engine rotation rate, current.
	VehicleEngineRpmCurrent float32 `json:"vehicle_engine_rpm_current"`

	// Throttle pedal after assists and overrides, 0 (off) to 1 (full).
	VehicleThrottle float32 `json:"vehicle_throttle"`

	// Brake pedal after assists and overrides, 0 (off) to 1 (full).
	VehicleBrake float32 `json:"vehicle_brake"`

	// Clutch pedal after assists and overrides, 0 (off) to 1 (full).
	VehicleClutch float32 `json:"vehicle_clutch"`

	// Steering after assists and overrides, -1 (full left) to 1 (full right).
	VehicleSteering float32 `json:"vehicle_steering"`

	// Handbrake after assists and overrides, 0 (off) to 1 (full).
	VehicleHandbrake float32 `json:"vehicle_handbrake"`

	// Time spent on current stage.
	StageCurrentTime float32 `json:"stage_current_time"`

	// Distance reached on current stage.
	StageCurrentDistance float64 `json:"stage_current_distance"`

	// Total length of current stage.
	StageLength float64 `json:"stage_length"`
}

// CurrentStageTime returns current stagetime as a formatted string. "03:42.583"
func (p *Packet) CurrentStageTime() (t string) {
	return Timespan(time.Duration(p.StageCurrentTime * float32(time.Second))).Format("04:05.000")
}

// InGameTime returns time spent ingame as a formatted string. "03:42.583"
func (p *Packet) InGameTime() (t string) {
	return Timespan(time.Duration(p.GameTotalTime * float32(time.Second))).Format("04:05.000")
}

// ToJSON returns the packet as marshaled JSON.
func (p *Packet) ToJSON() ([]byte, error) {
	return jettison.Marshal(p)
}

// StageProgress returns the current stage's progress as a percentage.
func (p *Packet) StageProgress() Percentage {
	return Percentage((p.StageCurrentDistance / p.StageLength) * 100)
}

// Throttle returns the current throttle value as a percentage.
func (p *Packet) Throttle() Percentage {
	return Percentage(p.VehicleThrottle * 100)
}

// Brake returns the current brake value as a percentage.
func (p *Packet) Brake() Percentage {
	return Percentage(p.VehicleBrake * 100)
}

// Clutch returns the current clutch value as a percentage. (100 = engaged)
func (p *Packet) Clutch() Percentage {
	return Percentage(p.VehicleClutch * 100)
}

// Handbrake returns the current handbrake value as a percentage.
func (p *Packet) Handbrake() Percentage {
	return Percentage(p.VehicleHandbrake * 100)
}

func (p *Packet) Kmph() float32 {
	return p.VehicleSpeed * MpsToKmph
}

func (p *Packet) Mph() float32 {
	return p.VehicleSpeed * MpsToMph
}

func NewPacket() *Packet {
	return &Packet{
		PacketUID:                 0,
		GameTotalTime:             0,
		GameDeltaTime:             0,
		GameFrameCount:            0,
		ShiftlightsFraction:       0,
		ShiftlightsRpmStart:       0,
		ShiftlightsRpmEnd:         0,
		ShiftlightsRpmValid:       false,
		VehicleGearIndex:          0,
		VehicleGearIndexNeutral:   0,
		VehicleGearIndexReverse:   0,
		VehicleGearMaximum:        0,
		VehicleSpeed:              0,
		VehicleTransmissionSpeed:  0,
		VehiclePositionX:          0,
		VehiclePositionY:          0,
		VehiclePositionZ:          0,
		VehicleVelocityX:          0,
		VehicleVelocityY:          0,
		VehicleVelocityZ:          0,
		VehicleAccelerationX:      0,
		VehicleAccelerationY:      0,
		VehicleAccelerationZ:      0,
		VehicleLeftDirectionX:     0,
		VehicleLeftDirectionY:     0,
		VehicleLeftDirectionZ:     0,
		VehicleForwardDirectionX:  0,
		VehicleForwardDirectionY:  0,
		VehicleForwardDirectionZ:  0,
		VehicleUpDirectionX:       0,
		VehicleUpDirectionY:       0,
		VehicleUpDirectionZ:       0,
		VehicleHubPositionBL:      0,
		VehicleHubPositionBR:      0,
		VehicleHubPositionFL:      0,
		VehicleHubPositionFR:      0,
		VehicleHubVelocityBL:      0,
		VehicleHubVelocityBR:      0,
		VehicleHubVelocityFL:      0,
		VehicleHubVelocityFR:      0,
		VehicleCpForwardSpeedBL:   0,
		VehicleCpForwardSpeedBR:   0,
		VehicleCpForwardSpeedFL:   0,
		VehicleCpForwardSpeedFR:   0,
		VehicleBrakeTemperatureBL: 0,
		VehicleBrakeTemperatureBR: 0,
		VehicleBrakeTemperatureFL: 0,
		VehicleBrakeTemperatureFR: 0,
		VehicleEngineRpmMax:       0,
		VehicleEngineRpmIdle:      0,
		VehicleEngineRpmCurrent:   0,
		VehicleThrottle:           0,
		VehicleBrake:              0,
		VehicleClutch:             0,
		VehicleSteering:           0,
		VehicleHandbrake:          0,
		StageCurrentTime:          0,
		StageCurrentDistance:      0,
		StageLength:               0,
	}
}
