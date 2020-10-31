import bpy
import sys
argv = sys.argv
argv = argv[argv.index("--") + 1:]

scn = bpy.data.scenes["Scene"]

index = int(argv[0])
xmin = float(argv[1])
ymin = float(argv[2])
xmax = float(argv[3])
ymax = float(argv[4])

# get the frame arg
frame = int(argv[5])

# go to the desired frame
scn.frame_set(frame)

# set a unique output path
# scn.render.filepath = "//render/" + str(index)
# for i in range(0, 4):
#     scn.render.filepath += "_" + argv[i]

scn.render.image_settings.file_format = 'OPEN_EXR'
scn.render.image_settings.color_mode = 'RGBA'
scn.render.image_settings.exr_codec = 'PIZ'
scn.render.image_settings.use_preview = False
scn.render.use_compositing = False
scn.render.use_sequencer = False
scn.render.use_save_buffers = False
scn.render.use_persistent_data = False

scn.cycles.use_cache = False
scn.cycles.debug_use_spatial_splits = False
scn.cycles.use_progressive_refine = False

percentage = max(1, min(10000, scn.render.resolution_percentage))
resx = int(scn.render.resolution_x * percentage / 100)
resy = int(scn.render.resolution_y * percentage / 100)

scn.render.use_border = True
scn.render.use_crop_to_border = True

scn.render.tile_x = max(4, min(64, (xmax - xmin + 1) // 8))
scn.render.tile_y = max(4, min(64, (ymax - ymin + 1) // 8))
print("Using tiles of size", scn.render.tile_x, scn.render.tile_y)
scn.render.threads_mode = 'AUTO'

# setup the render border
scn.render.border_min_x = (xmin + 0.5) / resx
scn.render.border_max_x = (xmax + 1.5) / resx
scn.render.border_min_y = (ymin + 0.5) / resy
scn.render.border_max_y = (ymax + 1.5) / resy

# render a still frame
bpy.ops.render.render(write_still=True)
print("Success.")