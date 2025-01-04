Encounter log version *15* documentation:

All lines begin with the time in MS since logging began and the line type.

_unitState_ refers to the following fields for a unit: unitId, health/max, magicka/max, stamina/max, ultimate/max, werewolf/max, shield, map NX, map NY, headingRadians.

_targetUnitState_ is replaced with an asterisk if the source and target are the same.

_equipmentInfo_ refers to the following fields for a piece of equipment: slot. id, isCP, level, trait, displayQuality, setId, enchantType, isEnchantCP, enchantLevel, enchantQuality.

__BEGIN_LOG__ - timeSinceEpocsMS, logVersion, realmName, language, gameVersion

__END_LOG__

__BEGIN_COMBAT__

__END_COMBAT__

__PLAYER_INFO__ - unitId, [longTermEffectAbilityId,...], [longTermEffectStackCounts,...], [_equipmentInfo_,...], [primaryAbilityId,...], [backupAbilityId,...]

__BEGIN_CAST__ - durationMS, channeled, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_

__END_CAST__ - endReason, castTrackId, interruptingAbilityId:optional, interruptingUnitId:optional

__COMBAT_EVENT__ - actionResult, damageType, powerType, hitValue, overflow, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_

__HEALTH_REGEN__ - effectiveRegen, _unitState_

__UNIT_ADDED__ - unitId, unitType, isLocalPlayer, playerPerSessionId, monsterId, isBoss, classId, raceId, name, displayName, characterId, level, championPoints, ownerUnitId, reaction, isGroupedWithLocalPlayer

__UNIT_CHANGED__ - unitId, classId, raceId, name, displayName, characterId, level, championPoints, ownerUnitId, reaction, isGroupedWithLocalPlayer

__UNIT_REMOVED__ - unitId

__EFFECT_CHANGED__ - changeType, stackCount, castTrackId, abilityId, _sourceUnitState_, _targetUnitState_, playerInitiatedRemoveCastTrackId:optional

__ABILITY_INFO__ - abilityId, name, iconPath, interruptible, blockable

__EFECT_INFO__ - abilityId, effectTye, statusEffectType, effectBarDisplayBehaviour, grantsSynergyAbilityId:optional

__MAP_INFO__ - id, name, texturePath

__ZONE_INFO__ - id, name, dungeonDifficulty

__TRIAL_INIT__ - id, inProgress, completed, startTimeMS, durationMS, success, finalScore

__BEGIN_TRIAL__ - id, startTimeMS

__END_TRIAL__ - id, durationMS, success, finalScore, finalVitalityBonus
